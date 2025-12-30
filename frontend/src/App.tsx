import { useState } from "react";
import "./App.css";
import { Sidebar } from "@/components/shared/Sidebar";
import { MainContent } from "@/components/shared/MainContent";

export default function App() {
	const [currentView, setCurrentView] = useState("models");

	return (
		<div className="flex h-screen overflow-hidden">
			<Sidebar currentView={currentView} onViewChange={setCurrentView} />
			<MainContent view={currentView} />
		</div>
	);
}
