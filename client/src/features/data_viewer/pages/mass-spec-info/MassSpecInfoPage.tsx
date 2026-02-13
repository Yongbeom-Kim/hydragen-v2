import { CircularProgress, Alert, Typography, Stack, Box } from "@mui/joy";
import { MassSpectrumDashboard } from "./components/mass-spectrum-dashboard";
import {
	useCompoundDetailQuery,
	useMassSpectrumQuery,
} from "../../hooks/query";

export function MassSpecInfoPageComponent({
	inchiHash,
}: {
	inchiHash: string;
}) {
	const { data, isLoading, error } = useCompoundDetailQuery({ inchiHash });

	if (isLoading) {
		return <CircularProgress sx={{ mt: 2 }} />;
	}

	if (!data) {
		return (
			<Alert color="danger" sx={{ mt: 2 }}>
				Failed to load compound.
				{error && (
					<Box component="span" sx={{ display: "block", mt: 1, fontSize: 13 }}>
						{String(error)}
					</Box>
				)}
			</Alert>
		);
	}

	return (
		<Box sx={{ maxWidth: 1100, mx: "auto", px: 2, py: 3 }}>
			<Stack direction="row" spacing={2} sx={{ alignItems: "center", mb: 2 }}>
				<Box
					component="img"
					src={data.imageUrl}
					alt={data.name}
					loading="lazy"
					sx={{
						width: 96,
						height: 96,
						borderRadius: "md",
						objectFit: "contain",
						bgcolor: "background.level1",
						border: "1px solid",
						borderColor: "divider",
						p: 0.5,
						flexShrink: 0,
					}}
				/>
				<Typography level="h1">
					Mass Spectra: {data?.name || data?.formula}
				</Typography>
			</Stack>
			<MassSpecInfoPageComponentInner inchiHash={inchiHash} />
		</Box>
	);
}

function MassSpecInfoPageComponentInner({ inchiHash }: { inchiHash: string }) {
	const { data, isLoading, error } = useMassSpectrumQuery({ inchiHash });

	if (isLoading) {
		return <CircularProgress sx={{ mt: 2 }} />;
	}

	if (!data || !data.items || data.items.length === 0) {
		return (
			<Alert color="danger" sx={{ mt: 2 }}>
				Failed to load mass spectrum.
				{error && (
					<Box component="span" sx={{ display: "block", mt: 1, fontSize: 13 }}>
						{String(error)}
					</Box>
				)}
			</Alert>
		);
	}

	return (
		<Stack>
			{data.items.map((spectrum, idx) => (
				<MassSpectrumDashboard key={spectrum.id || idx} spectrum={spectrum} />
			))}
		</Stack>
	);
}
