import { useQuery } from "@tanstack/react-query";
import {
	getCompoundDetail,
	getCompoundList,
	getMassSpectrum,
} from "../api/data-api";

type UseCompoundListQueryParams = {
	page: number;
	pageSize: number;
};

export const useCompoundListQuery = ({
	page,
	pageSize,
}: UseCompoundListQueryParams) => {
	return useQuery({
		queryKey: ["compounds", page, pageSize],
		queryFn: () => getCompoundList(page, pageSize),
	});
};

export const useCompoundDetailQuery = ({ inchiHash }: { inchiHash: string }) =>
	useQuery({
		queryKey: ["compound_detail", inchiHash],
		queryFn: () => getCompoundDetail(inchiHash),
	});

export const useMassSpectrumQuery = ({ inchiHash }: { inchiHash: string }) =>
	useQuery({
		queryKey: ["mass_spectrum_detail", inchiHash],
		queryFn: () => getMassSpectrum(inchiHash),
	});
