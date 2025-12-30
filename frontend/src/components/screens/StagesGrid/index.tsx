import { useState, useEffect } from "react";
import { GetStages } from "../../../../wailsjs/go/main/App";
import { main } from "../../../../wailsjs/go/models";
import { Button } from "@/components/ui/button";
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "@/components/ui/select";
import { ChevronLeft, ChevronRight, RefreshCw } from "lucide-react";
import { MMDContentCard } from "../../shared/MMDContentCard";

interface StagesGridProps {
	onShowDetail: (
		type: "model" | "stage",
		item: {
			id: string;
			name: string;
			screenshots: string[];
			description: string;
			originalPath: string;
		}
	) => void;
}

export function StagesGrid({ onShowDetail }: StagesGridProps) {
	const [paginatedData, setPaginatedData] =
		useState<main.PaginatedStages | null>(null);
	const [loading, setLoading] = useState(true);
	const [page, setPage] = useState(1);
	const [perPage, setPerPage] = useState(100);

	const loadStages = async () => {
		setLoading(true);
		try {
			const data = await GetStages(page, perPage);
			setPaginatedData(data);
		} catch (error) {
			console.error("Error loading stages:", error);
		} finally {
			setLoading(false);
		}
	};

	useEffect(() => {
		loadStages();
	}, [page, perPage]);

	const handlePageChange = (newPage: number) => {
		if (
			paginatedData &&
			newPage >= 1 &&
			newPage <= paginatedData.totalPages
		) {
			setPage(newPage);
		}
	};

	const handlePerPageChange = (value: string) => {
		setPerPage(Number.parseInt(value));
		setPage(1); // Reset to first page when changing page size
	};

	if (loading && !paginatedData) {
		return (
			<div className="flex items-center justify-center h-64">
				<div className="text-muted-foreground">Loading stages...</div>
			</div>
		);
	}

	if (!paginatedData || paginatedData.stages.length === 0) {
		return (
			<div className="flex items-center justify-center h-64">
				<div className="text-muted-foreground">No stages found</div>
			</div>
		);
	}

	return (
		<div className="space-y-6">
			{/* Controls */}
			<div className="flex items-center justify-between">
				<div className="flex items-center gap-4">
					<span className="text-sm text-muted-foreground">
						Showing {(page - 1) * perPage + 1} to{" "}
						{Math.min(page * perPage, paginatedData.total)} of{" "}
						{paginatedData.total} stages
					</span>
					<Select value={perPage.toString()} onValueChange={handlePerPageChange}>
						<SelectTrigger className="w-32">
							<SelectValue />
						</SelectTrigger>
						<SelectContent>
							<SelectItem value="100">100 per page</SelectItem>
							<SelectItem value="1000">1000 per page</SelectItem>
							<SelectItem value="10000">10000 per page</SelectItem>
						</SelectContent>
					</Select>
				</div>

				<Button
					variant="outline"
					size="sm"
					onClick={loadStages}
					disabled={loading}
				>
					<RefreshCw className={`w-4 h-4 mr-2 ${loading ? "animate-spin" : ""}`} />
					Refresh
				</Button>
			</div>

			{/* Stages Grid */}
			<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
				{paginatedData.stages.map((stage) => (
					<MMDContentCard
						key={stage.id}
						id={stage.id}
						name={stage.name}
						screenshots={stage.screenshots}
						description={stage.description}
						onClick={() => onShowDetail("stage", stage)}
					/>
				))}
			</div>

			{/* Pagination */}
			<div className="flex items-center justify-center gap-2">
				<Button
					variant="outline"
					size="sm"
					onClick={() => handlePageChange(page - 1)}
					disabled={page === 1 || loading}
				>
					<ChevronLeft className="w-4 h-4" />
					Previous
				</Button>

				<div className="flex items-center gap-2">
					{/* Show first page */}
					{page > 3 && (
						<>
							<Button
								variant="outline"
								size="sm"
								onClick={() => handlePageChange(1)}
								disabled={loading}
							>
								1
							</Button>
							{page > 4 && <span className="text-muted-foreground">...</span>}
						</>
					)}

					{/* Show pages around current page */}
					{Array.from({ length: 5 }, (_, i) => page - 2 + i)
						.filter((p) => p >= 1 && p <= paginatedData.totalPages)
						.map((p) => (
							<Button
								key={p}
								variant={p === page ? "default" : "outline"}
								size="sm"
								onClick={() => handlePageChange(p)}
								disabled={loading}
							>
								{p}
							</Button>
						))}

					{/* Show last page */}
					{page < paginatedData.totalPages - 2 && (
						<>
							{page < paginatedData.totalPages - 3 && (
								<span className="text-muted-foreground">...</span>
							)}
							<Button
								variant="outline"
								size="sm"
								onClick={() => handlePageChange(paginatedData.totalPages)}
								disabled={loading}
							>
								{paginatedData.totalPages}
							</Button>
						</>
					)}
				</div>

				<Button
					variant="outline"
					size="sm"
					onClick={() => handlePageChange(page + 1)}
					disabled={page === paginatedData.totalPages || loading}
				>
					Next
					<ChevronRight className="w-4 h-4" />
				</Button>
			</div>
		</div>
	);
}
