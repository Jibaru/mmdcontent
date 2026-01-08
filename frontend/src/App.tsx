import { useState } from "react";
import "./App.css";
import { Sidebar } from "@/components/shared/Sidebar";
import { MainContent } from "@/components/shared/MainContent";

export type ViewState = {
	view: string;
	detailType?: "model" | "stage" | "motion";
	detailItem?: {
		id: string;
		name: string;
		screenshots: string[];
		description: string;
		originalPath: string;
	};
};

export default function App() {
	const [viewState, setViewState] = useState<ViewState>({ view: "models" });

	const handleViewChange = (view: string) => {
		setViewState({ view });
	};

	const handleShowDetail = (
		type: "model" | "stage" | "motion",
		item: {
			id: string;
			name: string;
			screenshots: string[];
			description: string;
			originalPath: string;
		}
	) => {
		setViewState({
			view: "detail",
			detailType: type,
			detailItem: item,
		});
	};

	const handleBackFromDetail = () => {
		if (viewState.detailType === "model") {
			setViewState({ view: "models" });
		} else if (viewState.detailType === "stage") {
			setViewState({ view: "stages" });
		} else if (viewState.detailType === "motion") {
			setViewState({ view: "motions" });
		}
	};

	return (
		<div className="flex h-screen overflow-hidden">
			<Sidebar currentView={viewState.view} onViewChange={handleViewChange} />
			<MainContent
				viewState={viewState}
				onShowDetail={handleShowDetail}
				onBackFromDetail={handleBackFromDetail}
			/>
		</div>
	);
}
