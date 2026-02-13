import {
	Box,
	Card,
	CircularProgress,
	Link as JoyLink,
	Table,
	Typography,
	IconButton,
	Sheet,
	Select,
	Option,
} from "@mui/joy";
import { useNavigate } from "@tanstack/react-router";
import {
	createColumnHelper,
	flexRender,
	getCoreRowModel,
	useReactTable,
} from "@tanstack/react-table";
import { ChevronLeft, ChevronRight } from "lucide-react";
import type { CompoundListItem } from "../../../types/data";

const columnHelper = createColumnHelper<CompoundListItem>();

const columns = [
	columnHelper.accessor("name", {
		header: "Compound",
		cell: (info) => info.getValue(),
	}),
	columnHelper.accessor("formula", {
		header: "Formula",
	}),
	columnHelper.accessor("molecularWeight", {
		header: "Molecular Weight",
		cell: (info) => (info.getValue() ? info.getValue()?.toFixed(4) : "N/A"),
	}),
	columnHelper.accessor("hasMassSpectrum", {
		header: "Mass Spectrum",
		cell: (info) => (info.getValue() ? "✅" : "❌"),
	}),
	columnHelper.accessor("inchiKey", {
		header: "InChIKey",
	}),
];

type Props = {
	data: CompoundListItem[];
	loading: boolean;
	page: number;
	pageSize: number;
	total: number;
	onPageChange: (page: number) => void;
	onPageSizeChange: (pageSize: number) => void;
};

export function CompoundsTable({
	data,
	loading,
	page,
	pageSize,
	total,
	onPageChange,
	onPageSizeChange,
}: Props) {
	const navigate = useNavigate();

	const table = useReactTable({
		data,
		columns,
		getCoreRowModel: getCoreRowModel(),
	});

	const pageCount = Math.max(1, Math.ceil(total / pageSize));

	return (
		<Card size="lg" variant="outlined" sx={{ mt: 2, gap: 2 }}>
			<Typography level="h2">Compounds</Typography>
			<Typography level="body-sm" color="neutral">
				Sorted by molecular weight ascending.
			</Typography>

			<Sheet variant="soft" sx={{ borderRadius: "md", overflow: "auto" }}>
				<Table stickyHeader hoverRow>
					<thead>
						{table.getHeaderGroups().map((headerGroup) => (
							<tr key={headerGroup.id}>
								{headerGroup.headers.map((header) => (
									<th key={header.id}>
										{flexRender(
											header.column.columnDef.header,
											header.getContext(),
										)}
									</th>
								))}
							</tr>
						))}
					</thead>
					<tbody>
						{loading ? (
							<tr>
								<td colSpan={columns.length}>
									<Box
										sx={{ display: "flex", justifyContent: "center", py: 3 }}
									>
										<CircularProgress />
									</Box>
								</td>
							</tr>
						) : (
							table.getRowModel().rows.map((row) => (
								<tr
									key={row.id}
									onClick={() =>
										navigate({
											to: "/data/compounds/$inchiHash",
											params: { inchiHash: row.original.inchiKey },
										})
									}
									style={{ cursor: "pointer" }}
								>
									{row.getVisibleCells().map((cell, index) => (
										<td key={cell.id}>
											{index === 0 ? (
												<JoyLink
													sx={{ textDecoration: "none", fontWeight: 700 }}
												>
													{flexRender(
														cell.column.columnDef.cell,
														cell.getContext(),
													)}
												</JoyLink>
											) : (
												flexRender(
													cell.column.columnDef.cell,
													cell.getContext(),
												)
											)}
										</td>
									))}
								</tr>
							))
						)}
					</tbody>
				</Table>
			</Sheet>

			<Box
				sx={{
					display: "flex",
					justifyContent: "space-between",
					alignItems: "center",
				}}
			>
				<Box sx={{ display: "flex", gap: 1, alignItems: "center" }}>
					<IconButton
						variant="outlined"
						size="sm"
						disabled={page <= 1}
						onClick={() => onPageChange(page - 1)}
					>
						<ChevronLeft size={18} />
					</IconButton>
					<Typography>
						Page {page} / {pageCount}
					</Typography>
					<IconButton
						variant="outlined"
						size="sm"
						disabled={page >= pageCount}
						onClick={() => onPageChange(page + 1)}
					>
						<ChevronRight size={18} />
					</IconButton>
				</Box>

				<Select
					size="sm"
					value={pageSize}
					onChange={(_, value) => value && onPageSizeChange(value)}
				>
					<Option value={10}>10 / page</Option>
					<Option value={20}>20 / page</Option>
					<Option value={50}>50 / page</Option>
				</Select>
			</Box>
		</Card>
	);
}
