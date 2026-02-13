import {
	Card,
	Chip,
	Divider,
	Link as JoyLink,
	Stack,
	Typography,
} from "@mui/joy";
import { Link } from "@tanstack/react-router";
import { useCompoundDetailQuery } from "@/features/data_viewer/hooks/query";
import { ChevronRight } from "lucide-react";
type Props = {
	inchiHash: string;
};

export function CompoundInfo({ inchiHash }: Props) {
	const { data } = useCompoundDetailQuery({ inchiHash });

	if (!data) {
		throw new Error(
			"Unreachable code (this component should not be mounted if data is null).",
		);
	}

	const { name, inchiKey, hasMassSpectrum } = data;
	return (
		<Card size="lg" variant="outlined" sx={{ mt: 2, gap: 1.5 }}>
			<Typography level="h2">{name}</Typography>
			<Typography level="body-sm">InChIKey: {inchiKey}</Typography>
			<Typography>
				<strong>Formula:</strong> {data.formula}
			</Typography>
			<Typography>
				<strong>SMILES:</strong> {data.smiles}
			</Typography>
			<Typography>
				<strong>InChI:</strong> {data.inchi}
			</Typography>

			<Divider />

			<Stack direction="row" spacing={1}>
				<Typography>
					<strong>Mass Spectrum:</strong>
				</Typography>
				{hasMassSpectrum ? (
					<MassSpecAvailableChip inchiKey={inchiKey} />
				) : (
					<MassSpecUnavailableChip />
				)}
			</Stack>
		</Card>
	);
}

function MassSpecAvailableChip({ inchiKey }: { inchiKey: string }) {
	return (
		<JoyLink
			component={Link}
			to={`/data/mass_spectra/${inchiKey}`}
			sx={{ textDecoration: "none" }}
		>
			<Chip color="success" startDecorator="✅">
				<div className="flex flex-row items-center gap-1">
					available <ChevronRight size="14" />
				</div>
			</Chip>
		</JoyLink>
	);
}

function MassSpecUnavailableChip() {
	return (
		<Chip color="danger" startDecorator="❌">
			Not available
		</Chip>
	);
}
