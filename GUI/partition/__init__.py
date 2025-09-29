# Ce fichier permet de transformer le répertoire en package Python
# Il permet d'importer les modules qu'il contient
# Vous pouvez également définir ici quels sous-modules seront exposés

# Exposer les classes publiques pour faciliter l'import
from .PartitionModel import Partition
from .PartitionService import PartitionService
from .PartitionPresenter import PartitionPresenter
from .PartitionPanel import PartitionPanel
from .PartitionWidget import PartitionWidget

__all__ = [
    "Partition",
    "PartitionService",
    "PartitionPresenter",
    "PartitionPanel",
    "PartitionWidget",
]
