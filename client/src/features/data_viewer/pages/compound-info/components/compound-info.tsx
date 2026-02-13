import {
	Box,
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

	const { name, inchiKey, hasMassSpectrum, imageUrl } = data;
	return (
		<Card size="lg" variant="outlined" sx={{ mt: 2, gap: 1.5 }}>
			<Stack direction="row" spacing={2} sx={{ alignItems: "center" }}>
				<Box
					component="img"
					src={imageUrl}
					alt={name}
					loading="lazy"
					sx={{
						width: 112,
						height: 112,
						borderRadius: "md",
						objectFit: "contain",
						bgcolor: "background.level1",
						border: "1px solid",
						borderColor: "divider",
						p: 0.5,
						flexShrink: 0,
					}}
				/>
				<Typography level="h2">{name}</Typography>
			</Stack>
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
