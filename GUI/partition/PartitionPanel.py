from PyQt6.QtWidgets import QFrame, QVBoxLayout, QLabel, QScrollArea, QWidget
from PyQt6.QtCore import Qt, pyqtSignal

from .PartitionWidget import PartitionWidget


class PartitionPanel(QFrame):
    """
    Panneau qui affiche la liste des partitions disponibles.
    """

    partition_selected = pyqtSignal(str)

    def __init__(self, parent=None):
        super().__init__(parent)

        self.setStyleSheet(
            """
            background-color: #0c0c23;
            border-radius: 15px;
            padding: 15px;
        """
        )

        self.init_ui()

    def init_ui(self):
        self.layout = QVBoxLayout(self)

        self.title_label = QLabel("Partitions disponibles")
        self.title_label.setStyleSheet(
            "color: white; font-size: 16px; font-weight: bold;"
        )
        self.title_label.setAlignment(Qt.AlignmentFlag.AlignCenter)
        self.layout.addWidget(self.title_label)

        self.info_label = QLabel(
            "Sélectionnez un disque et cliquez sur 'Charger les partitions'"
        )
        self.info_label.setStyleSheet("color: white; font-size: 14px;")
        self.info_label.setAlignment(Qt.AlignmentFlag.AlignCenter)
        self.layout.addWidget(self.info_label)

        self.scroll_area = QScrollArea()
        self.scroll_area.setWidgetResizable(True)
        self.scroll_area.setStyleSheet(
            """
            QScrollArea {
                border: none;
                background-color: transparent;
            }
            QScrollBar:vertical {
                background-color: #1e1e2e;
                width: 10px;
                border-radius: 5px;
            }
            QScrollBar::handle:vertical {
                background-color: #3D4455;
                min-height: 20px;
                border-radius: 5px;
            }
            QScrollBar::add-line:vertical, QScrollBar::sub-line:vertical {
                height: 0px;
            }
        """
        )

        self.partitions_container = QWidget()
        self.partitions_layout = QVBoxLayout(self.partitions_container)
        self.partitions_layout.setContentsMargins(0, 0, 0, 0)
        self.partitions_layout.setSpacing(10)

        self.partitions_layout.addStretch()

        self.scroll_area.setWidget(self.partitions_container)
        self.layout.addWidget(self.scroll_area)

    def show_loading(self, disk_name):
        self.clear_partitions()
        self.info_label.setText(f"Chargement des partitions de {disk_name}...")

    def show_error(self, message):
        self.clear_partitions()
        self.info_label.setText(message)
        self.info_label.setStyleSheet("color: #ff5555; font-size: 14px;")

    def show_empty(self, message):
        self.clear_partitions()
        self.info_label.setText(message)
        self.info_label.setStyleSheet("color: white; font-size: 14px;")

    def show_partitions(self, partitions):
        self.clear_partitions()

        if not partitions:
            self.show_empty("Aucune partition trouvée")
            return

        self.info_label.setText(f"Partitions disponibles ({len(partitions)})")
        self.info_label.setStyleSheet("color: white; font-size: 14px;")

        for partition in partitions:
            partition_widget = PartitionWidget(partition)
            partition_widget.load_files_requested.connect(self.partition_selected)

            self.partitions_layout.insertWidget(
                self.partitions_layout.count() - 1, partition_widget
            )

    def clear_partitions(self):
        while self.partitions_layout.count() > 1:
            item = self.partitions_layout.takeAt(0)
            if item.widget():
                item.widget().deleteLater()
