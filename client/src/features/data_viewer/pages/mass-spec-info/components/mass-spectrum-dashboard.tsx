import {
	Accordion,
	AccordionDetails,
	AccordionGroup,
	AccordionSummary,
	Box,
	Divider,
	Typography,
} from "@mui/joy";
import ReactECharts from "echarts-for-react";
import type { MassSpectrumItem } from "../../../types/data";
import React from "react";

type Props = {
	spectrum: MassSpectrumItem;
};

export function MassSpectrumDashboard({ spectrum }: Props) {
	return (
		<Box sx={{ mt: 2, p: 0 }}>
			<MassSpectrumMetadataAccordion spectrum={spectrum} />
			<MassSpectrumGraph spectrum={spectrum} />
		</Box>
	);
}

function MassSpectrumMetadataAccordion({ spectrum }: Props) {
	const metadataRows = [
		{
			label: "Molecular Weight:",
			value: spectrum.molecularWeight,
		},
		{
			label: "Exact Mass:",
			value: spectrum.exactMass ?? "N/A",
		},
		{
			label: "Precursor m/z:",
			value: spectrum.precursorMz ?? "N/A",
		},
		{
			label: "Precursor Type:",
			value: spectrum.precursorType ?? "N/A",
		},
		{
			label: "Ion Mode:",
			value: spectrum.ionMode ?? "N/A",
		},
		{
			label: "Collision Energy:",
			value: spectrum.collisionEnergy ?? "N/A",
		},
		{
			label: "Instrument:",
			value: spectrum.instrument ?? "N/A",
		},
		{
			label: "Instrument Type:",
			value: spectrum.instrumentType ?? "N/A",
		},
		{
			label: "SPLASH:",
			value: spectrum.splash ?? "N/A",
		},
		{
			label: "DB Number:",
			value: spectrum.dbNumber,
		},
	];

	return (
		<AccordionGroup
			variant="outlined"
			sx={{
				borderRadius: "lg",
				borderBottomRightRadius: 0,
				borderBottomLeftRadius: 0,
				borderBottom: "none",
			}}
		>
			<Accordion defaultExpanded={false}>
				<AccordionSummary
					sx={{
						borderRadius: "lg",
						"&>button": {
							borderRadius: "lg",
							borderBottomRightRadius: 0,
							borderBottomLeftRadius: 0,
						},
					}}
				>
					<Typography level="h2" sx={{ p: 2 }}>
						Record #{spectrum.id} · {spectrum.source}
					</Typography>
				</AccordionSummary>
				<AccordionDetails>
					<Box
						sx={{
							display: "grid",
							gridTemplateColumns: "max-content auto",
							rowGap: 2,
							columnGap: 6,
							m: 2,
						}}
						component="dl"
					>
						{metadataRows.map((row, idx) => (
							<React.Fragment key={row.label}>
								<Typography component="dt" fontWeight="bold" level="body-lg">
									{row.label}
								</Typography>
								<Typography component="dd" level="body-lg">
									{row.value}
								</Typography>
							</React.Fragment>
						))}
					</Box>
					<Divider sx={{ mt: 2 }} />
				</AccordionDetails>
			</Accordion>
		</AccordionGroup>
	);
}

function MassSpectrumGraph({ spectrum }: Props) {
	const dataPoints = spectrum.mZ.map((mz, index) => [
		mz,
		spectrum.peaks[index] ?? 0,
	]);
	return (
		<Box
			sx={{
				py: 2,
				px: 4,
				borderLeft: "1px solid",
				borderRight: "1px solid",
				borderBottom: "1px solid",
				borderColor: "divider",
				borderTop: "none",
				borderRadius: "lg",
				borderTopLeftRadius: 0,
				borderTopRightRadius: 0,
			}}
		>
			<Typography level="h3" sx={{ mb: 1 }}>
				Spectrum visualisation
			</Typography>
			<ReactECharts
				style={{ height: 420 }}
				option={{
					backgroundColor: "#fff",
					title: {
						text: "Mass Spectrum",
						left: "center",
					},
					tooltip: { trigger: "axis" },
					xAxis: {
						name: "m/z",
						nameLocation: "middle",
						nameGap: 30,
						type: "value",
					},
					yAxis: {
						name: "Relative Intensity",
						type: "value",
					},
					series: [
						{
							type: "bar",
							data: dataPoints,
							barMaxWidth: 8,
							itemStyle: {
								color: "#5470c6",
							},
						},
					],
					grid: { top: 60, left: 60, right: 20, bottom: 60 },
				}}
			/>
		</Box>
	);
}
