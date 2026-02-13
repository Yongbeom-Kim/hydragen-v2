import { createFileRoute } from "@tanstack/react-router";
import { CompoundsPage } from "@/features/data_viewer/pages/compounds-page/CompoundsPage";

export const Route = createFileRoute("/data/")({
	validateSearch: (search: Record<string, unknown>) => ({
		page: Number(search.page ?? 1),
		pageSize: Number(search.pageSize ?? 20),
	}),
	component: CompoundsRoute,
});

function CompoundsRoute() {
	const { page, pageSize } = Route.useSearch();
	return <CompoundsPage page={page} pageSize={pageSize} />;
}
