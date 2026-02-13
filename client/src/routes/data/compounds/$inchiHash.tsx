import { createFileRoute } from "@tanstack/react-router";
import { CompoundInfoPage as CompoundInfoPageComponent } from "@/features/data_viewer/pages/compound-info/CompoundInfoPage";

export const Route = createFileRoute("/data/compounds/$inchiHash")({
	component: CompoundInfoPage,
});

function CompoundInfoPage() {
	const { inchiHash } = Route.useParams();
	return <CompoundInfoPageComponent inchiHash={inchiHash} />;
}
