import re

from PyQt6.QtWidgets import (
    QFrame,
    QVBoxLayout,
    QHBoxLayout,
    QLabel,
    QPushButton,
    QProgressBar,
    QDialog,
)
from PyQt6.QtCore import Qt, pyqtSignal


class PartitionWidget(QFrame):
    """
    Widget qui affiche les informations d'une partition.
    """

    load_files_requested = pyqtSignal(str)

    def __init__(self, partition, parent=None):
        super().__init__(parent)
        self.partition = partition

        self.setStyleSheet(
            """
            background-color: #1e1e2e;
            border-radius: 8px;
            margin: 5px;
            padding: 10px;
        """
        )

        self.init_ui()

    def init_ui(self):
        layout = QVBoxLayout(self)

        header_layout = QHBoxLayout()

        name_label = QLabel(self.partition.name)
        name_label.setStyleSheet("color: white; font-size: 14px; font-weight: bold;")
        header_layout.addWidget(name_label)

        load_files_btn = QPushButton("Charger")
        load_files_btn.setObjectName(f"load_files_{self.partition.name}")
        load_files_btn.setStyleSheet(
            """
            QPushButton {
                background-color: #3D4455;
                color: white;
                border-radius: 5px;
                padding: 4px 8px;
                font-size: 12px;
            }
            QPushButton:hover {
                background-color: #66A3FF;
            }
            QPushButton:pressed {
                background-color: #4D5466;
            }
        """
        )
        load_files_btn.clicked.connect(self.on_load_button_clicked)
        header_layout.addWidget(load_files_btn)

        details_btn = QPushButton("Détails")
        load_files_btn.setObjectName(f"details_partition_{self.partition.name}")
        details_btn.setStyleSheet(
            """
            QPushButton {
                background-color: #3D4455;
                color: white;
                border-radius: 5px;
                padding: 4px 8px;
                font-size: 12px;
            }
            QPushButton:hover {
                background-color: #66A3FF;
            }
            QPushButton:pressed {
                background-color: #4D5466;
            }
        """
        )

        details_btn.clicked.connect(self.details_button_clicked)

        header_layout.addWidget(details_btn)

        layout.addLayout(header_layout)

        self.progress_bar = QProgressBar()
        self.progress_bar.setStyleSheet(
            """
            QProgressBar {
                border-radius: 5px;
                background-color: #3D4455;
                color: white;
                font-size: 12px;
            }
            QProgressBar::chunk {
                background-color: #66A3FF;
                border-radius: 5px;
            }
            """
        )
        self.progress_bar.setTextVisible(True)

        self.update_progress_bar()

        layout.addWidget(self.progress_bar)

        label_space = QVBoxLayout()

        label_label = QLabel(f"Label : {self.partition.label}")
        label_label.setStyleSheet("color: white; font-size: 13px; padding-top: 5px;")
        label_space.addWidget(label_label)

        fs_label = QLabel(f"Système de fichiers : {self.partition.filesystem}")
        fs_label.setStyleSheet("color: white; font-size: 13px; padding-top: 5px;")
        label_space.addWidget(fs_label)

        layout.addLayout(label_space)

    def on_load_button_clicked(self):
        self.load_files_requested.emit(self.partition.name)

    def details_button_clicked(self):
        dialog = QDialog(self)
        dialog.setWindowTitle(f"Détails de {self.partition.name}")

        dialog.setStyleSheet(
            """
            QDialog {
                background-color: #2E2E3E;
                color: white;
                padding: 15px;
                border-radius: 15px;  /* Rounded corners */
                border: 2px solid #44475a; /* Border color */
            }
            """
        )

        dialog.setFixedSize(600, 600)
        dialog.setWindowFlags(Qt.WindowType.FramelessWindowHint | Qt.WindowType.Dialog)

        parent_geometry = self.geometry()
        parent_center_x = parent_geometry.x() + parent_geometry.width() // 2
        parent_center_y = parent_geometry.y() + parent_geometry.height() // 2

        dialog.move(parent_center_x - 200, parent_center_y - 200)

        layout = QVBoxLayout(dialog)

        details = [
            f"Nom : {self.partition.name}",
            f"Label : {self.partition.label}",
            f"Système de fichiers : {self.partition.filesystem}",
            f"Taille totale : {self.partition.size}",
            f"Point de montage : {self.partition.mount_point}",
            f"Taille fs : {self.partition.fs_size}",
            f"Utilisé : {self.partition.fs_used}",
            f"uuid : {self.partition.uuid}",
            f"Taille secteur physique : {self.partition.physical_sector_size}",
            f"Taille secteur logique : {self.partition.logical_sector_size}",
        ]

        for detail in details:
            label = QLabel(detail)
            label.setStyleSheet("font-size: 14px; padding: 5px; color: white")
            layout.addWidget(label)

        close_btn = QPushButton("Fermer")
        close_btn.setStyleSheet(
            """
            QPushButton {
                background-color: #66A3FF;
                color: white;
                border-radius: 5px;
                padding: 6px 12px;
                font-size: 14px;
            }
            QPushButton:hover {
                background-color: #4D90FE;
            }
            QPushButton:pressed {
                background-color: #3D7DD4;
            }
            """
        )
        close_btn.clicked.connect(dialog.close)
        layout.addWidget(close_btn, alignment=Qt.AlignmentFlag.AlignCenter)

        dialog.setLayout(layout)
        dialog.exec()

    def update_progress_bar(self):
        if self.partition.size != "0":

            total_size = self.convert_size_to_bytes(self.partition.size)
            used_space = self.convert_size_to_bytes(self.partition.fs_used)

            used_percentage = (used_space / total_size) * 100

            self.progress_bar.setValue(int(used_percentage))

            self.progress_bar.setFormat(f"     Utilisé: {int(used_percentage)}%")
        else:
            self.progress_bar.setValue(0)
            self.progress_bar.setFormat("Aucune donnée")

    def convert_size_to_bytes(self, size_str):
        size_str = size_str.strip().upper()
        size_map = {"B": 1, "K": 1024, "M": 1024**2, "G": 1024**3, "T": 1024**4}

        match = re.match(r"([\d\.]+)([BKMGT])", size_str)
        if match:
            value, unit = match.groups()
            return float(value) * size_map[unit]
        else:
            return 0
