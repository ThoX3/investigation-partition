from PyQt6.QtWidgets import (
    QWidget,
    QVBoxLayout,
    QTableWidget,
    QTableWidgetItem,
    QHeaderView,
    QLabel,
    QPushButton,
    QHBoxLayout,
    QWidget,
    QDialog,
    QScrollArea,
    QGridLayout,
    QTextEdit,
    QLineEdit,
)
from PyQt6.QtCore import Qt, pyqtSignal

from .FileService import FileService

from utils import *
from PyQt6.QtGui import QIcon
import os


class FileTableWidget(QWidget):
    """Widget personnalisé pour afficher une liste de fichiers."""

    folder_opened = pyqtSignal(str)

    def __init__(self, parent=None):
        super().__init__(parent)
        self.current_partition_id = None
        self.current_path = ""

        self.layout = QVBoxLayout(self)
        self.setLayout(self.layout)

        self.file_type_filter = QLineEdit()
        self.file_type_filter.setPlaceholderText(
            "Saisir une extension (ex: .txt, .png)"
        )
        self.file_type_filter.textChanged.connect(self.apply_filter)
        self.file_type_filter.setStyleSheet("color: white;")
        self.layout.addWidget(self.file_type_filter)

        self.message_label = QLabel("")
        self.message_label.setAlignment(Qt.AlignmentFlag.AlignCenter)
        self.message_label.setStyleSheet("color: white; font-size: 14px;")
        self.layout.addWidget(self.message_label)

        self.table = QTableWidget(0, 2)
        self.table.setHorizontalHeaderLabels(["Nom", "Actions"])
        self.table.horizontalHeader().setSectionResizeMode(
            QHeaderView.ResizeMode.Stretch
        )
        self.table.setColumnWidth(0, 300)
        self.table.setColumnWidth(1, 200)

        self.table.setVerticalScrollBarPolicy(Qt.ScrollBarPolicy.ScrollBarAlwaysOff)
        self.table.horizontalHeader().setVisible(False)
        self.table.verticalHeader().setVisible(False)
        self.table.setShowGrid(False)
        self.table.setEditTriggers(QTableWidget.EditTrigger.NoEditTriggers)
        self.table.cellDoubleClicked.connect(self.on_item_double_clicked)
        self.apply_style()
        self.layout.addWidget(self.table)

        self.display_message("Sélectionnez une partition pour voir les fichiers")

    def apply_style(self):
        self.table.setStyleSheet(
            """
            QTableWidget {
                background-color: #1e1e2e;
                color: white;
                border: none;
                gridline-color: transparent;
            }
            QHeaderView::section {
                background-color: #444;
                color: white;
                padding: 6px;
                font-weight: bold;
                border: none;
            }
            QTableWidget::item {
                border-bottom: none;
                padding: 6px;
            }
            QTableWidget::item:selected {
                background-color: #3D4455;
            }
            QPushButton {
                background-color: #444;
                color: white;
                border-radius: 5px;
                padding: 4px 6px;
                min-width: 90px;
                font-size: 12px;
                border: 1px solid #888;  
            }
            QPushButton:hover {
                background-color: #666;
            }
            QPushButton:pressed {
                background-color: #555;
                border: 1px solid #bbb;  /* Even lighter border when pressed */
            }
        """
        )

        self.table.setStyleSheet(
            self.table.styleSheet()
            + """
            QWidget {
                background-color: #1e1e2e;
            }
        """
        )

    def display_message(self, message):
        self.message_label.setText(message)
        self.message_label.show()
        self.table.hide()

    def load_partition_files(self, partition_id, path=""):
        self.current_partition_id = partition_id
        self.current_path = path

        files = FileService.fetch_files(partition_id, path)

        self.message_label.hide()
        self.table.show()

        if not files.get("files"):
            self.table.setRowCount(1)
            back_item = QTableWidgetItem("⬅ Revenir en arrière")
            back_item.setFlags(Qt.ItemFlag.ItemIsSelectable | Qt.ItemFlag.ItemIsEnabled)
            self.table.setItem(0, 0, back_item)
            return

        selected_filter = self.file_type_filter.text().strip()
        if selected_filter:
            files["files"] = [
                file
                for file in files["files"]
                if file.get("name", "").endswith(selected_filter)
            ]

        self.table.setRowCount(len(files["files"]) + 1)

        back_item = QTableWidgetItem("⬅ Revenir en arrière")
        back_item.setFlags(Qt.ItemFlag.ItemIsSelectable | Qt.ItemFlag.ItemIsEnabled)
        self.table.setItem(0, 0, back_item)

        for row, file in enumerate(files.get("files")):
            name = file.get("name", "Inconnu")
            extension = file.get("type", "")

            name_item = QTableWidgetItem(name)
            if extension == "directory":
                name_item.setIcon(QIcon("icons/open-folder.png"))
                name_item.setFlags(
                    Qt.ItemFlag.ItemIsSelectable | Qt.ItemFlag.ItemIsEnabled
                )
            else:
                if extension == "file":
                    ext = os.path.splitext(name)[1].lower().lstrip(".")
                    icon_path = f"icons/{ext}.png"
                    if not os.path.exists(icon_path):
                        icon_path = "icons/file.png"

                    name_item.setIcon(QIcon(icon_path))
                else:
                    name_item.setIcon(QIcon("icons/link.png"))
                name_item.setFlags(Qt.ItemFlag.NoItemFlags)
            self.table.setItem(row + 1, 0, name_item)

            btn_widget = QWidget()
            btn_layout = QHBoxLayout(btn_widget)
            btn_layout.setContentsMargins(0, 0, 0, 0)
            btn_layout.setSpacing(10)

            details_btn = QPushButton("Détails")
            details_btn.clicked.connect(lambda _, f=file: self.show_details(f))
            btn_layout.addWidget(details_btn)

            if extension == "file":
                show_blocks_btn = QPushButton("Blocks")
                show_blocks_btn.clicked.connect(
                    lambda _, f=file: self.show_blocks_dialog(f)
                )
                btn_layout.addWidget(show_blocks_btn)

            btn_widget.setLayout(btn_layout)
            self.table.setCellWidget(row + 1, 1, btn_widget)

    def on_item_double_clicked(self, row, column):
        if column != 0:
            return

        if row == 0:
            self.go_back()
            return

        name_item = self.table.item(row, 0)
        if name_item:
            folder_name = name_item.text()
            self.current_path = (
                f"{self.current_path}/{folder_name}"
                if self.current_path
                else folder_name
            )
            self.load_partition_files(self.current_partition_id, self.current_path)
            self.folder_opened.emit(self.current_path)

    def show_details(self, file):
        dialog = QDialog(self)
        dialog.setWindowTitle(f"Détails de {file.get('name')}")

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
            f"Nom : {file.get('name', 'inconnu')}",
            f"Type : {file.get('type', 'inconnu')}",
            f"Symlink : {file.get('is_symlink', 'inconnu')}",
            f"Taille (bytes) : {file.get('size_bytes', 'inconnu')}",
            f"Permissions : {file.get('permissions', 'inconnu')}",
            f"Liens durs : {file.get('hard_links', 'inconnu')}",
            f"Inode : {file.get('inode', 'inconnu')}",
            f"UID propriétaire : {file.get('owner_uid', 'inconnu')}",
            f"GID propriétaire : {file.get('owner_gid', 'inconnu')}",
            f"Taille de bloc : {file.get('block_size', 'inconnu')}",
            f"Blocs alloués : {file.get('blocks_allocated', 'inconnu')}",
            f"Dernière modification : {file.get('last_modified', 'inconnu')}",
            f"Dernier accès : {file.get('last_access', 'inconnu')}",
        ]

        for detail in details:
            label = QLabel(detail)
            label.setStyleSheet("font-size: 14px; padding: 5px; color : white;")
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

    def go_back(self):
        self.current_path = "/".join(self.current_path.split("/")[:-1]) or ""
        self.load_partition_files(self.current_partition_id, self.current_path)
        self.folder_opened.emit(self.current_path)

    def show_blocks_dialog(self, file):
        path = f"{self.current_path}/{file['name']}"
        current_partition = self.current_partition_id

        blocks = getBlocksOfFile(current_partition, path)

        dialog = QDialog(self)
        dialog.setWindowTitle("Blocks Viewer")
        dialog.setMinimumSize(600, 400)
        dialog.setStyleSheet("background-color: #1e1e2e; color: white;")

        main_layout = QVBoxLayout(dialog)

        title_label = QLabel(f"Blocks of the File {file.get('name', '')} :")
        title_label.setAlignment(Qt.AlignmentFlag.AlignCenter)
        title_label.setStyleSheet(
            "font-size: 16px; font-weight: bold; margin-bottom: 10px;"
        )
        main_layout.addWidget(title_label)

        scroll_area = QScrollArea()
        scroll_area.setWidgetResizable(True)
        scroll_area.setStyleSheet("border: none;")

        container = QWidget()
        grid_layout = QGridLayout(container)
        grid_layout.setSpacing(10)

        row, col = 0, 0
        max_columns = 5
        for block_id in sorted(blocks.keys(), key=int):
            block_btn = QPushButton(f"Block {block_id}")
            block_btn.setFixedSize(80, 80)
            block_btn.setStyleSheet(
                """
                QPushButton {
                    background-color: #444;
                    color: white;
                    border: 1px solid #888;
                    border-radius: 10px;
                    font-size: 12px;
                }
                QPushButton:hover {
                    background-color: #666;
                }
            """
            )

            block_btn.clicked.connect(
                lambda _, b=block_id: self.show_block_content(self, blocks, b)
            )
            grid_layout.addWidget(block_btn, row, col)

            col += 1
            if col >= max_columns:
                col = 0
                row += 1

        container.setLayout(grid_layout)
        scroll_area.setWidget(container)
        main_layout.addWidget(scroll_area)

        dialog.setLayout(main_layout)
        dialog.exec()

    def show_block_content(self, parent, blocks, block_id):
        block_content = blocks.get(block_id, "No Data Available")

        content_dialog = QDialog(parent)
        content_dialog.setWindowTitle(f"Block {block_id} Content")
        content_dialog.setMinimumSize(500, 300)
        content_dialog.setStyleSheet("background-color: #1e1e2e; color: white;")

        layout = QVBoxLayout(content_dialog)

        label = QLabel(f"Content of Block {block_id}:")
        label.setStyleSheet("font-size: 14px; font-weight: bold;")
        layout.addWidget(label)

        text_area = QTextEdit()
        text_area.setText(block_content)
        text_area.setReadOnly(True)
        text_area.setStyleSheet("background-color: #333; color: white; border: none;")
        layout.addWidget(text_area)

        content_dialog.setLayout(layout)
        content_dialog.exec()

    def apply_filter(self):
        if self.current_partition_id:
            self.load_partition_files(self.current_partition_id, self.current_path)
