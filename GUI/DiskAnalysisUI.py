from PyQt6.QtWidgets import (
    QApplication,
    QWidget,
    QVBoxLayout,
    QLabel,
    QHBoxLayout,
    QFrame,
    QTabWidget,
    QPushButton,
)
from PyQt6.QtCore import Qt
import sys

from ColoredDisk import *
from DiskService import *
from file_service.FileTableWidget import FileTableWidget

from partition.PartitionPanel import PartitionPanel
from partition.PartitionPresenter import PartitionPresenter


class DiskAnalysisUI(QWidget):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("Analyse du disque")
        self.setGeometry(100, 100, 1000, 600)
        self.setStyleSheet("background-color: #000000;")

        self.current_disk = None
        self.disks = []
        self.tab_widget = None

        layout = QHBoxLayout()
        layout.addLayout(self.create_left_panel())
        layout.addWidget(self.create_right_panel())

        self.setLayout(layout)

        self.partition_presenter = PartitionPresenter(self)

        self.setup_connections()

    def create_left_panel(self):
        left_panel = QVBoxLayout()

        left_panel.addWidget(self.create_disk_tabs())
        left_panel.addWidget(self.create_file_table())

        return left_panel

    def create_disk_tabs(self):
        self.tab_widget = QTabWidget()
        self.tab_widget.setObjectName("disk_tab_widget")
        self.tab_widget.setStyleSheet(
            """
            background-color: #0c0c23;
            border-radius: 15px;
            padding: 15px;
            color: white;
                                 
            QTabWidget::pane {
                border: 1px solid #C2C7CB;
                border-radius: 4px;
                color: white;
            }
            QTabBar::tab {
                background-color: #A3C1DA;
                padding: 10px;
                border: 1px solid #C2C7CB;
                border-radius: 4px;
                color: white;
            }
            QTabBar::tab:selected {
                background-color: #66A3FF;
                color: white;
            }
            QTabBar::tab:hover {
                background-color: #B8D0FF;
                color: white;
            }
        """
        )

        disks, error_msg = fetch_disks()
        self.disks = disks or []

        if error_msg:
            error_tab = QWidget()
            layout = QVBoxLayout(error_tab)
            label = QLabel(f"Erreur lors de la récupération des disques:\n{error_msg}")
            label.setAlignment(Qt.AlignmentFlag.AlignCenter)
            label.setStyleSheet("color: red; font-size: 14px;")
            layout.addWidget(label)
            self.tab_widget.addTab(error_tab, "Erreur API")
            return self.tab_widget

        if not self.disks:
            empty_tab = QWidget()
            layout = QVBoxLayout(empty_tab)
            label = QLabel("Aucun disque détecté")
            label.setAlignment(Qt.AlignmentFlag.AlignCenter)
            label.setStyleSheet("color: white; font-size: 14px;")
            layout.addWidget(label)
            self.tab_widget.addTab(empty_tab, "Aucun disque")
            return self.tab_widget

        for disk in self.disks:
            tab = QWidget()
            layout = QVBoxLayout(tab)

            colored_disk = ColoredDisk(disk)
            colored_disk.setFixedSize(260, 260)
            layout.addWidget(colored_disk, alignment=Qt.AlignmentFlag.AlignCenter)

            legend_layout = QHBoxLayout()
            layout.addLayout(legend_layout)

            legend_layout.addStretch(1)
            self.add_legend_item(legend_layout, "#33ff57", "Free Space")
            self.add_legend_item(legend_layout, "#ff5733", "Lost Space")
            self.add_legend_item(legend_layout, "#3357ff", "Used Space")
            legend_layout.addStretch(1)

            disk_name = QLabel(f"Nom : {disk.get('name', 'Inconnu')}")
            disk_capacity = QLabel(f"Capacité : {disk.get('capacity_gb', '0')} Go")

            disk_name.setStyleSheet("font-size: 14px; color: white;")
            disk_capacity.setStyleSheet("font-size: 14px; color: white;")

            layout.addWidget(disk_name, alignment=Qt.AlignmentFlag.AlignCenter)
            layout.addWidget(disk_capacity, alignment=Qt.AlignmentFlag.AlignCenter)

            load_button = QPushButton("Charger les partitions")
            load_button.setObjectName(f"load_button_{disk.get('name', 'unknown')}")
            load_button.setStyleSheet(
                """
                QPushButton {
                    background-color: #3D4455;
                    color: white;
                    border-radius: 5px;
                    padding: 8px 16px;
                    font-size: 14px;
                }
                QPushButton:hover {
                    background-color: #66A3FF;
                }
                QPushButton:pressed {
                    background-color: #4D5466;
                }
            """
            )
            layout.addWidget(load_button, alignment=Qt.AlignmentFlag.AlignCenter)

            tab_name = disk.get("name", "Disque ?")
            if "/" in tab_name:
                tab_name = tab_name.split("/")[-1]

            self.tab_widget.addTab(tab, tab_name)

        return self.tab_widget

    def create_file_table(self):
        panel_table_container = QFrame()
        panel_table_container.setStyleSheet(
            """
            background-color: #0c0c23;
            border-radius: 15px;
            padding: 15px;
        """
        )

        layout = QVBoxLayout(panel_table_container)

        title_label = QLabel("Fichiers de la partition")
        title_label.setStyleSheet("color: white; font-size: 14px; font-weight: bold;")
        title_label.setAlignment(Qt.AlignmentFlag.AlignCenter)
        layout.addWidget(title_label)

        self.file_table = FileTableWidget()
        layout.addWidget(self.file_table)

        return panel_table_container

    def create_right_panel(self):
        self.partition_panel = PartitionPanel()
        return self.partition_panel

    def setup_connections(self):
        if not self.tab_widget:
            return

        self.tab_widget.currentChanged.connect(self.on_tab_changed)

        for i in range(self.tab_widget.count()):
            tab = self.tab_widget.widget(i)
            for child in tab.findChildren(QPushButton):
                if child.objectName().startswith("load_button_"):
                    disk_name = child.objectName().replace("load_button_", "")
                    child.clicked.connect(
                        lambda checked, name=disk_name: self.on_load_partitions_clicked(
                            name
                        )
                    )

        if hasattr(self, "partition_panel"):
            self.partition_panel.partition_selected.connect(self.on_partition_selected)

    def on_tab_changed(self, index):
        if index >= 0 and index < len(self.disks):
            disk_name = self.disks[index].get("name", None)
            self.partition_presenter.set_current_disk(disk_name)

    def on_load_partitions_clicked(self, disk_name):
        self.partition_presenter.load_partitions(disk_name)

    def on_partition_selected(self, partition_name):
        self.file_table.load_partition_files(partition_name)

    def show_partition_loading(self, disk_name):
        self.partition_panel.show_loading(disk_name)

    def show_partition_error(self, message):
        self.partition_panel.show_error(message)

    def show_partition_empty(self, message):
        self.partition_panel.show_empty(message)

    def show_partitions(self, partitions):
        self.partition_panel.show_partitions(partitions)

    def load_partition_files(self, partition_name):
        self.file_table.load_partition_files(partition_name)

    def add_legend_item(self, layout, color, text):
        label = QLabel()
        label.setFixedSize(20, 20)
        label.setStyleSheet(f"background-color: {color}; border: 1px solid black;")
        text_label = QLabel(text)
        text_label.setStyleSheet("margin-left: 5px;")

        item_layout = QHBoxLayout()
        item_layout.addWidget(label)
        item_layout.addWidget(text_label)
        layout.addLayout(item_layout)


if __name__ == "__main__":
    app = QApplication(sys.argv)
    window = DiskAnalysisUI()
    window.show()
    sys.exit(app.exec())
