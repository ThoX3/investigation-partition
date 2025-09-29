from PyQt6.QtWidgets import QWidget
from PyQt6.QtGui import QPainter, QBrush, QColor
from PyQt6.QtCore import Qt
from utils import *


class ColoredDisk(QWidget):
    def __init__(self, disk, parent=None):
        super().__init__(parent)
        self.disk = disk
        self.capacity = disk.get("capacity_gb", "0")
        self.free_space = self.calculate_spaces()[0]
        self.lost_space = self.calculate_spaces()[1]
        self.used_space = self.calculate_spaces()[2]

    def paintEvent(self, event):
        painter = QPainter(self)
        painter.setRenderHint(QPainter.RenderHint.Antialiasing)

        total = self.capacity if self.capacity > 0 else 1

        segments = [
            ("#33ff57", int((self.free_space / total) * 360)),
            ("#ff5733", int((self.lost_space / total) * 360)),
            ("#3357ff", int((self.used_space / total) * 360)),
        ]

        start_angle = 0
        radius = 250
        center_x, center_y = 37, 37

        for color, angle in segments:
            painter.setPen(Qt.PenStyle.NoPen)
            painter.setBrush(QBrush(QColor(color)))
            painter.drawPie(10, 10, radius, radius, start_angle * 16, angle * 16)
            start_angle += angle

        painter.setBrush(QBrush(QColor("#0c0c23")))
        painter.drawEllipse(center_x, center_y, 200, 200)

    def calculate_spaces(self):
        partitions = get_partitions_from_disk(self.disk.get("name"))

        free_space = 0.0
        lost_space = 0.0
        used_space = 0.0

        for partition in partitions:
            size = self.convert_to_gb(partition.get("size", "0"))
            fs_size = self.convert_to_gb(partition.get("fs_size", "0"))
            fs_used = self.convert_to_gb(partition.get("fs_used", "0"))

            used_space += fs_used
            free_space += fs_size - fs_used
            lost_space += size - fs_size

        return free_space, lost_space, used_space
    
    def convert_to_gb(self, value):
        if isinstance(value, str):
            value = value.strip().lower()
            if any(unit in value for unit in ["to", "t", "tb"]):
                return float(value.replace("to", "").replace("t", "").replace("tb", "").strip()) * 1024
            elif any(unit in value for unit in ["go", "g", "gb"]):
                return float(value.replace("go", "").replace("g", "").replace("gb", "").strip())
            elif any(unit in value for unit in ["mo", "m", "mb"]):
                return float(value.replace("mo", "").replace("m", "").replace("mb", "").strip()) / 1024
        return float(value)
