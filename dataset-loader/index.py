from loaders.massbank_data_loader import MassbankDataLoader

loaders = [
	MassbankDataLoader(
		'MassBankDataLoader',
		'https://github.com/MassBank/MassBank-data/releases/download/2025.10/MassBank_NISTformat.msp'
	)
]

for loader in loaders:
	loader.load()