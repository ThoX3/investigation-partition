from .PartitionService import PartitionService


class PartitionPresenter:
    """
    Présentateur pour gérer la logique concernant les partitions.
    Fait le lien entre la vue et le modèle/service.
    """

    def __init__(self, view):
        self.view = view
        self.current_disk = None
        self.partitions = []

    def set_current_disk(self, disk_name):
        self.current_disk = disk_name
        self.partitions = []

    def load_partitions(self, disk_name=None):
        if not disk_name:
            disk_name = self.current_disk

        if not disk_name:
            self.view.show_partition_error("Aucun disque sélectionné")
            return

        self.view.show_partition_loading(disk_name)

        partitions, error = PartitionService.get_partitions_for_disk(disk_name)

        if error:
            self.view.show_partition_error(
                f"Erreur lors du chargement des partitions : {error}"
            )
            return

        if not partitions:
            self.view.show_partition_empty(f"Aucune partition trouvée pour {disk_name}")
            return

        self.partitions = partitions
        self.view.show_partitions(partitions)

    def load_partition_files(self, partition_name):
        self.view.load_partition_files(partition_name)
