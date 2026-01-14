import { useState, useEffect } from "react";
import { GetStages, SearchStages } from "../../../../wailsjs/go/handlers/Stages";
import { entities } from "../../../../wailsjs/go/models";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "@/components/ui/select";
import { ChevronLeft, ChevronRight, RefreshCw, Search, X } from "lucide-react";
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
		useState<entities.Pagination_MMDContent_internal_entities_Stage_ | null>(null);
	const [searchResults, setSearchResults] = useState<entities.Stage[] | null>(null);
	const [loading, setLoading] = useState(true);
	const [searching, setSearching] = useState(false);
	const [page, setPage] = useState(1);
	const [perPage, setPerPage] = useState(10);
	const [searchQuery, setSearchQuery] = useState("");

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

	const handleSearch = async () => {
		if (!searchQuery.trim()) {
			setSearchResults(null);
			return;
		}

		setSearching(true);
		try {
			const results = await SearchStages(searchQuery, 1000); // Limit to top 1000 results
			setSearchResults(results);
		} catch (error) {
			console.error("Error searching stages:", error);
		} finally {
			setSearching(false);
		}
	};

	const handleClearSearch = () => {
		setSearchQuery("");
		setSearchResults(null);
	};

	useEffect(() => {
		if (!searchResults) {
			loadStages();
		}
	}, [page, perPage, searchResults]);

	useEffect(() => {
		// Debounce search
		const timer = setTimeout(() => {
			if (searchQuery.trim()) {
				handleSearch();
			} else {
				setSearchResults(null);
			}
		}, 500);

		return () => clearTimeout(timer);
	}, [searchQuery]);

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

	// Use search results if searching, otherwise use paginated data
	const displayData = searchResults || paginatedData?.data || [];
	const isSearching = searchResults !== null;

	if (loading && !paginatedData && !searchResults) {
		return (
			<div className="flex items-center justify-center h-64">
				<div className="text-muted-foreground">Loading stages...</div>
			</div>
		);
	}

	if (!isSearching && (!paginatedData || paginatedData.data.length === 0)) {
		return (
			<div className="flex items-center justify-center h-64">
				<div className="text-muted-foreground">No stages found</div>
			</div>
		);
	}

	return (
		<div className="space-y-6">
			{/* Search Bar */}
			<div className="flex items-center gap-4">
				<div className="relative flex-1 max-w-xl">
					<Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
					<Input
						placeholder="Search stages by description... (powered by AI)"
						value={searchQuery}
						onChange={(e) => setSearchQuery(e.target.value)}
						className="pl-10 pr-10"
					/>
					{searchQuery && (
						<Button
							variant="ghost"
							size="sm"
							className="absolute right-1 top-1/2 -translate-y-1/2 h-7 w-7 p-0"
							onClick={handleClearSearch}
						>
							<X className="h-4 w-4" />
						</Button>
					)}
				</div>
				{isSearching && (
					<div className="text-sm text-muted-foreground">
						{searchResults.length} result{searchResults.length !== 1 ? 's' : ''} found
					</div>
				)}
			</div>

			{/* Controls */}
			{!isSearching && (
				<div className="flex items-center justify-between">
					<div className="flex items-center gap-4">
						<span className="text-sm text-muted-foreground">
							Showing {(page - 1) * perPage + 1} to{" "}
							{Math.min(page * perPage, paginatedData?.total || 0)} of{" "}
							{paginatedData?.total || 0} stages
						</span>
						<Select value={perPage.toString()} onValueChange={handlePerPageChange}>
							<SelectTrigger className="w-32">
								<SelectValue />
							</SelectTrigger>
							<SelectContent>
								<SelectItem value="5">5 per page</SelectItem>
								<SelectItem value="10">10 per page</SelectItem>
								<SelectItem value="50">50 per page</SelectItem>
								<SelectItem value="100">100 per page</SelectItem>
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
			)}

			{/* Stages Grid */}
			{searching ? (
				<div className="flex items-center justify-center h-64">
					<div className="flex flex-col items-center gap-2">
						<div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900" />
						<span className="text-muted-foreground">Searching with AI...</span>
					</div>
				</div>
			) : displayData.length === 0 ? (
				<div className="flex items-center justify-center h-64">
					<div className="text-muted-foreground">
						{isSearching ? "No results found for your search" : "No stages found"}
					</div>
				</div>
			) : (
				<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
					{displayData.map((stage) => (
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
			)}

			{/* Pagination - only show when not searching */}
			{!isSearching && paginatedData && (
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
			)}
		</div>
	);
}
