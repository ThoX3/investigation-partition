class Partition:
    """
    Modèle pour représenter une partition de disque.
    Encapsule les données d'une partition avec les propriétés nécessaires.
    """

    def __init__(self, data=None):
        data = data or {}
        self.name = data.get("name", "Inconnu")
        self.filesystem = data.get("fs_type", "Inconnu")
        self.mount_point = data.get("mount_point", "")

        self.size = data.get("size", "0")
        self.fs_used = data.get("fs_used", "0")
        self.fs_size = data.get("fs_size", "0")
        self.uuid = data.get("uuid", "")
        self.label = data.get("part_type_name", "Sans étiquette")

        self.physical_sector_size = data.get("physical_sector_size", 0)
        self.logical_sector_size = data.get("logical_sector_size", 0)

    def __str__(self):
        return f"Partition {self.name} ({self.size_gb} {self.size_unity}, {self.filesystem})"

    @property
    def formatted_size(self):
        self.size_unity += "o"
        return f"{self.size_gb:.2f} {self.size_unity}"
