import CssBaseline from "@mui/joy/CssBaseline";
import { CssVarsProvider, extendTheme } from "@mui/joy/styles";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

const theme = extendTheme({
	colorSchemes: {
		light: {
			palette: {
				primary: {
					solidBg: "#0b6bcb",
				},
			},
		},
	},
});

export function getContext() {
	const queryClient = new QueryClient();
	return {
		queryClient,
	};
}

export function Provider({
	children,
	queryClient,
}: {
	children: React.ReactNode;
	queryClient: QueryClient;
}) {
	return (
		<QueryClientProvider client={queryClient}>
			<CssVarsProvider theme={theme} defaultMode="light">
				<CssBaseline />
				{children}
			</CssVarsProvider>
		</QueryClientProvider>
	);
}
