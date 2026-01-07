import { useState, useEffect } from "react";
import { GetImageAsBase64 } from "../../../../wailsjs/go/handlers/Images";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Dialog, DialogContent } from "@/components/ui/dialog";
import { ArrowLeft, Copy, FolderOpen, ZoomIn, ZoomOut, X } from "lucide-react";
import { BrowserOpenURL } from "../../../../wailsjs/runtime/runtime";

interface MMDContentDetailProps {
	type: "model" | "stage";
	item: {
		id: string;
		name: string;
		screenshots: string[];
		description: string;
		originalPath: string;
	};
	onBack: () => void;
}

export function MMDContentDetail({ type, item, onBack }: MMDContentDetailProps) {
	const [images, setImages] = useState<Map<number, string>>(new Map());
	const [selectedImageIndex, setSelectedImageIndex] = useState<number | null>(null);
	const [zoomLevel, setZoomLevel] = useState(100);
	const [loadingImages, setLoadingImages] = useState(true);

	// Load all images
	useEffect(() => {
		const loadImages = async () => {
			setLoadingImages(true);
			const imageMap = new Map<number, string>();

			for (let i = 0; i < item.screenshots.length; i++) {
				try {
					const base64 = await GetImageAsBase64(item.screenshots[i]);
					imageMap.set(i, base64);
				} catch (error) {
					console.error(`Error loading image ${i}:`, error);
				}
			}

			setImages(imageMap);
			setLoadingImages(false);
		};

		loadImages();
	}, [item.screenshots]);

	const handleCopyPath = async () => {
		try {
			await navigator.clipboard.writeText(item.originalPath);
		} catch (error) {
			console.error("Error copying path:", error);
		}
	};

	const handleCopyFolder = async () => {
		try {
			// Extract folder path from original path
			const folderPath = item.originalPath.substring(0, item.originalPath.lastIndexOf("\\"));
			await navigator.clipboard.writeText(folderPath);
		} catch (error) {
			console.error("Error copying folder path:", error);
		}
	};

	const handleOpenFolder = () => {
		// Extract folder path and open it
		const folderPath = item.originalPath.substring(0, item.originalPath.lastIndexOf("\\"));
		BrowserOpenURL(`file:///${folderPath}`);
	};

	const handleImageClick = (index: number) => {
		setSelectedImageIndex(index);
		setZoomLevel(100);
	};

	const handleCloseModal = () => {
		setSelectedImageIndex(null);
		setZoomLevel(100);
	};

	const handleZoomIn = () => {
		setZoomLevel(prev => Math.min(prev + 25, 300));
	};

	const handleZoomOut = () => {
		setZoomLevel(prev => Math.max(prev - 25, 50));
	};

	return (
		<div className="space-y-6">
			{/* Header with back button */}
			<div className="flex items-center gap-4">
				<Button variant="outline" size="sm" onClick={onBack}>
					<ArrowLeft className="w-4 h-4 mr-2" />
					Back to {type === "model" ? "Models" : "Stages"}
				</Button>
			</div>

			{/* Main info card */}
			<Card>
				<CardHeader>
					<div className="flex items-start justify-between">
						<div>
							<CardTitle className="text-2xl">{item.name}</CardTitle>
							<CardDescription className="mt-2">ID: {item.id}</CardDescription>
						</div>
						<div className="flex gap-2">
							<Button variant="outline" size="sm" onClick={handleCopyPath}>
								<Copy className="w-4 h-4 mr-2" />
								Copy Path
							</Button>
							<Button variant="outline" size="sm" onClick={handleCopyFolder}>
								<Copy className="w-4 h-4 mr-2" />
								Copy Folder
							</Button>
							<Button variant="outline" size="sm" onClick={handleOpenFolder}>
								<FolderOpen className="w-4 h-4 mr-2" />
								Open Folder
							</Button>
						</div>
					</div>
				</CardHeader>
				<CardContent className="space-y-6">
					{/* Description */}
					<div>
						<h3 className="text-sm font-semibold mb-2">Description</h3>
						<p className="text-sm text-muted-foreground whitespace-pre-wrap">
							{item.description || "No description available"}
						</p>
					</div>

					{/* Original Path */}
					<div>
						<h3 className="text-sm font-semibold mb-2">Original Path</h3>
						<code className="text-xs bg-muted px-2 py-1 rounded block break-all">
							{item.originalPath}
						</code>
					</div>
				</CardContent>
			</Card>

			{/* Screenshots card */}
			<Card>
				<CardHeader>
					<CardTitle>Screenshots ({item.screenshots.length})</CardTitle>
				</CardHeader>
				<CardContent>
					{loadingImages ? (
						<div className="flex items-center justify-center h-64">
							<div className="text-muted-foreground">Loading images...</div>
						</div>
					) : (
						<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
							{item.screenshots.map((_, index) => (
								<div
									key={index}
									className="relative aspect-video bg-muted rounded-lg overflow-hidden cursor-pointer hover:ring-2 hover:ring-primary transition-all"
									onClick={() => handleImageClick(index)}
								>
									{images.get(index) ? (
										<img
											src={images.get(index)}
											alt={`Screenshot ${index + 1}`}
											className="w-full h-full object-cover"
										/>
									) : (
										<div className="w-full h-full flex items-center justify-center">
											<span className="text-muted-foreground text-sm">
												Loading...
											</span>
										</div>
									)}
									<div className="absolute bottom-2 right-2 bg-black/50 text-white text-xs px-2 py-1 rounded">
										{index + 1}
									</div>
								</div>
							))}
						</div>
					)}
				</CardContent>
			</Card>

			{/* Image zoom modal */}
			<Dialog open={selectedImageIndex !== null} onOpenChange={handleCloseModal}>
				<DialogContent className="max-w-7xl w-full h-[90vh] p-0">
					<div className="relative w-full h-full bg-black/95 flex flex-col">
						{/* Modal header */}
						<div className="flex items-center justify-between p-4 bg-black/50">
							<span className="text-white font-medium">
								Screenshot {selectedImageIndex !== null ? selectedImageIndex + 1 : 0} of {item.screenshots.length}
							</span>
							<div className="flex items-center gap-2">
								<Button
									variant="ghost"
									size="sm"
									onClick={handleZoomOut}
									className="text-white hover:text-white hover:bg-white/20"
								>
									<ZoomOut className="w-4 h-4" />
								</Button>
								<span className="text-white text-sm w-16 text-center">
									{zoomLevel}%
								</span>
								<Button
									variant="ghost"
									size="sm"
									onClick={handleZoomIn}
									className="text-white hover:text-white hover:bg-white/20"
								>
									<ZoomIn className="w-4 h-4" />
								</Button>
								<Button
									variant="ghost"
									size="sm"
									onClick={handleCloseModal}
									className="text-white hover:text-white hover:bg-white/20"
								>
									<X className="w-4 h-4" />
								</Button>
							</div>
						</div>

						{/* Modal content */}
						<div className="flex-1 overflow-auto flex items-center justify-center p-4">
							{selectedImageIndex !== null && images.get(selectedImageIndex) && (
								<img
									src={images.get(selectedImageIndex)}
									alt={`Screenshot ${selectedImageIndex + 1}`}
									style={{
										width: `${zoomLevel}%`,
										maxWidth: "none",
									}}
									className="object-contain"
								/>
							)}
						</div>
					</div>
				</DialogContent>
			</Dialog>
		</div>
	);
}
