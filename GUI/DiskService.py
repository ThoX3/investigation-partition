import requests

API_URL_DISK = "http://localhost:8080/disks"


def fetch_disks():
    try:
        response = requests.get(API_URL_DISK, timeout=5)
        response.raise_for_status()
        disks = response.json()

        if not isinstance(disks, list):
            return None, "Format de réponse invalide"

        return disks, None
    except requests.exceptions.RequestException as e:
        return None, str(e)
