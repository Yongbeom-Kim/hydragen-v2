import { Alert, Box, Typography } from "@mui/joy";
import { useNavigate } from "@tanstack/react-router";
import { CompoundsTable } from "@/features/data_viewer/pages/compounds-page/components/CompoundsTable";
import { useCompoundListQuery } from "../../hooks/query";

export function CompoundsPage({
	page,
	pageSize,
}: {
	page: number;
	pageSize: number;
}) {
	return (
		<Box sx={{ maxWidth: 1200, mx: "auto", px: 2, py: 3 }}>
			<Typography level="h1">Compounds Data</Typography>
			<CompoundsPageData page={page} pageSize={pageSize} />
		</Box>
	);
}

function CompoundsPageData({
	page,
	pageSize,
}: {
	page: number;
	pageSize: number;
}) {
	const { data, isPending } = useCompoundListQuery({ page, pageSize });
	const navigate = useNavigate();

	if (!isPending && !data) {
		return (
			<Alert color="danger" sx={{ mt: 2 }}>
				Failed to load compounds list.
			</Alert>
		);
	}

	const {
		items = [],
		page: currentPage = page,
		pageSize: currentPageSize = pageSize,
		total = 0,
	} = data || {};

	return (
		<CompoundsTable
			data={items}
			loading={isPending}
			page={currentPage}
			pageSize={currentPageSize}
			total={total}
			onPageChange={(nextPage) =>
				navigate({
					to: "/data",
					search: { page: nextPage, pageSize },
				})
			}
			onPageSizeChange={(nextPageSize) =>
				navigate({
					to: "/data",
					search: { page: 1, pageSize: nextPageSize },
				})
			}
		/>
	);
}
