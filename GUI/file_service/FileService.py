import requests

class FileService:
    """Service pour communiquer avec l'API de fichiers."""
    
    BASE_URL = "http://localhost:8080"
    
    @staticmethod
    def fetch_files(partition_id, path=""):
        try:
            url = f"{FileService.BASE_URL}/files?partition={partition_id}&path={path}"
            response = requests.get(url, timeout=5)
            response.raise_for_status()
            files = response.json()
                        
            return files
        except requests.exceptions.RequestException as e:
            return str(e)