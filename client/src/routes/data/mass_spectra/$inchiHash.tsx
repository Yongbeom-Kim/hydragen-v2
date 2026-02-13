import { Alert, Box, CircularProgress, Stack, Typography } from "@mui/joy";
import { createFileRoute } from "@tanstack/react-router";
import { MassSpectrumDashboard } from "@/features/data_viewer/pages/mass-spec-info/components/mass-spectrum-dashboard";
import {
	useCompoundDetailQuery,
	useMassSpectrumQuery,
} from "@/features/data_viewer/hooks/query";
import { MassSpecInfoPageComponent } from "@/features/data_viewer/pages/mass-spec-info/MassSpecInfoPage";

export const Route = createFileRoute("/data/mass_spectra/$inchiHash")({
	component: MassSpectrumPage,
});

function MassSpectrumPage() {
	const { inchiHash } = Route.useParams();
	return <MassSpecInfoPageComponent inchiHash={inchiHash} />;
}
