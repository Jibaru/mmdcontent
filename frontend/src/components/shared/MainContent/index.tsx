import { ModelsGrid } from "../../screens/ModelsGrid";
import { StagesGrid } from "../../screens/StagesGrid";
import { MotionsGrid } from "../../screens/MotionsGrid";
import { MMDContentDetail } from "../../screens/MMDContentDetail";
import type { ViewState } from "../../../App";

interface MainContentProps {
	viewState: ViewState;
	onShowDetail: (
		type: "model" | "stage" | "motion",
		item: {
			id: string;
			name: string;
			screenshots: string[];
			description: string;
			originalPath: string;
		}
	) => void;
	onBackFromDetail: () => void;
}

export function MainContent({ viewState, onShowDetail, onBackFromDetail }: MainContentProps) {
	const { view, detailType, detailItem } = viewState;

	// Get page title based on view
	const getPageTitle = () => {
		if (view === "detail") {
			if (detailType === "model") return "Model Details";
			if (detailType === "stage") return "Stage Details";
			if (detailType === "motion") return "Motion Details";
		}
		if (view === "models") return "Models Library";
		if (view === "stages") return "Stages Library";
		if (view === "motions") return "Motions Library";
		return "Main Dashboard";
	};

	// Get breadcrumb based on view
	const getBreadcrumb = () => {
		if (view === "detail") {
			if (detailType === "model") return "Models / Details";
			if (detailType === "stage") return "Stages / Details";
			if (detailType === "motion") return "Motions / Details";
		}
		if (view === "models") return "Models";
		if (view === "stages") return "Stages";
		if (view === "motions") return "Motions";
		return "Dashboard";
	};

	return (
		<div className="flex-1 overflow-auto bg-gray-50">
			{/* Header */}
			<div className="bg-white border-b border-border px-8 py-4">
				<div className="flex items-center gap-2 text-sm text-muted-foreground mb-2">
					<span>Pages</span>
					<span>/</span>
					<span>{getBreadcrumb()}</span>
				</div>
				<h1 className="text-3xl font-bold">{getPageTitle()}</h1>
			</div>

			{/* Content */}
			<div className="p-8">
				{view === "models" && <ModelsGrid onShowDetail={onShowDetail} />}
				{view === "stages" && <StagesGrid onShowDetail={onShowDetail} />}
				{view === "motions" && <MotionsGrid onShowDetail={onShowDetail} />}
				{view === "detail" && detailType && detailItem && (
					<MMDContentDetail
						type={detailType}
						item={detailItem}
						onBack={onBackFromDetail}
					/>
				)}
				{view === "dashboard" && (
					<div className="text-center text-muted-foreground py-12">
						Dashboard statistics coming soon...
					</div>
				)}
			</div>
		</div>
	);
}
