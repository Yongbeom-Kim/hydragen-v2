import { Typography, CircularProgress, Alert, Box } from "@mui/joy";
import { CompoundInfo } from "./components/compound-info";
import { useCompoundDetailQuery } from "../../hooks/query";

export function CompoundInfoPage({ inchiHash }: { inchiHash: string }) {
	return (
		<Box sx={{ maxWidth: 1000, mx: "auto", px: 2, py: 3 }}>
			<Typography level="h1">Compound Information</Typography>
			<CompoundInfoPageInner inchiHash={inchiHash} />
		</Box>
	);
}

function CompoundInfoPageInner({ inchiHash }: { inchiHash: string }) {
	const { isPending, data } = useCompoundDetailQuery({ inchiHash });

	if (isPending) {
		return <CircularProgress sx={{ mt: 2 }} />;
	}

	if (!data) {
		return (
			<Alert color="danger" sx={{ mt: 2 }}>
				Failed to load compound.
			</Alert>
		);
	}

	return <CompoundInfo inchiHash={inchiHash} />;
}
