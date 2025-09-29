import requests

SERVER_IP = "localhost"
SERVER_PORT = 8080
PROTOCOL = "http"
BASE_URL = f"{PROTOCOL}://{SERVER_IP}:{SERVER_PORT}"
DISKS_ENDPOINT = "/disks"
PARTITIONS_ENDPOINT = "/partitions"
FILES_ENDPOINT = "/files"
BLOCKS_ENDPOINT = "/blocks"


def get_disks():
    url = f"{BASE_URL}{DISKS_ENDPOINT}"
    try:
        response = requests.get(url)
        response.raise_for_status()
        disks_data = response.json()

        disks_list = [(disk["name"], disk["capacity_gb"]) for disk in disks_data]
        return disks_list
    except requests.exceptions.RequestException as e:
        print(f"Erreur lors de la récupération des disques: {e}")
        return []


def get_partitions_from_disk(disk_name):
    url = f"{BASE_URL}{PARTITIONS_ENDPOINT}?disk={disk_name}"
    try:
        response = requests.get(url)
        response.raise_for_status()
        partitions_data = response.json()

        partitions_dict = [dict(partition) for partition in partitions_data]
        return partitions_dict
    except requests.exceptions.RequestException as e:
        print(
            f"Erreur lors de la récupération des partitions du disque {disk_name}: {e}"
        )
        return []


def getFilesAndfolderList(partitionName, path=None, filter=None):
    params = {"partition": partitionName}
    if path:
        params["path"] = path
    if filter:
        params["filter"] = filter

    url = f"{BASE_URL}{FILES_ENDPOINT}"

    try:
        response = requests.get(url, params=params)
        response.raise_for_status()

        return response.json()
    except requests.exceptions.RequestException as e:
        print(f"Erreur lors de la récupération des fichiers: {e}")
        return []


def getBlocksOfFile(partition, path):
    params = {"partition": partition, "path": path}
    url = f"{BASE_URL}{BLOCKS_ENDPOINT}"

    try:
        response = requests.get(url, params=params)
        response.raise_for_status()

        return response.json()
    except requests.exceptions.RequestException as e:
        print(f"Erreur lors de la récupération des blocs du fichier: {e}")
        return []
