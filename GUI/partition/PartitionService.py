import requests
from .PartitionModel import Partition


class PartitionService:
    """
    Service pour communiquer avec l'API concernant les partitions.
    Traite les requêtes et transforme les réponses en objets métier.
    """

    BASE_URL = "http://localhost:8080"
    PARTITIONS_ENDPOINT = "/partitions"

    @staticmethod
    def get_partitions_for_disk(disk_name):
        try:
            url = f"{PartitionService.BASE_URL}{PartitionService.PARTITIONS_ENDPOINT}"
            params = {"disk": disk_name}
            response = requests.get(url, params=params, timeout=5)
            response.raise_for_status()
            partitions_data = response.json()

            if not isinstance(partitions_data, list):
                return None, "Format de réponse invalide"

            partitions = []
            for partition_data in partitions_data:
                partition_data["disk_name"] = disk_name
                partitions.append(Partition(partition_data))

            return partitions, None
        except requests.exceptions.RequestException as e:
            return None, str(e)

    @staticmethod
    def get_partition_details(partition_name):
        try:
            url = f"{PartitionService.BASE_URL}{PartitionService.PARTITIONS_ENDPOINT}"
            params = {"name": partition_name}
            response = requests.get(url, params=params, timeout=5)
            response.raise_for_status()
            partition_data = response.json()

            if not partition_data:
                return None, f"Partition {partition_name} non trouvée"

            return Partition(partition_data[0]), None
        except requests.exceptions.RequestException as e:
            return None, str(e)
