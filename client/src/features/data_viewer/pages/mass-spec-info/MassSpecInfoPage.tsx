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
			<Typography level="h1">
				Mass Spectra: {data?.name || data?.formula}
			</Typography>
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
