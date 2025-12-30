import { useState, useEffect } from "react";
import { GetImageAsBase64 } from "../../../../wailsjs/go/main/App";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { ChevronLeft, ChevronRight } from "lucide-react";

interface MMDContentCardProps {
	id: string;
	name: string;
	screenshots: string[];
	description: string;
	onClick?: () => void;
}

export function MMDContentCard({
	id,
	name,
	screenshots,
	description,
	onClick,
}: MMDContentCardProps) {
	const [currentImageIndex, setCurrentImageIndex] = useState(0);
	const [loadedImages, setLoadedImages] = useState<Map<number, string>>(new Map());
	const [loading, setLoading] = useState(true);
	const [loadingImage, setLoadingImage] = useState(false);

	// Load the current image
	useEffect(() => {
		const loadCurrentImage = async () => {
			if (screenshots.length === 0) {
				setLoading(false);
				return;
			}

			// Check if image is already loaded
			if (loadedImages.has(currentImageIndex)) {
				setLoading(false);
				return;
			}

			setLoadingImage(true);
			try {
				const base64Image = await GetImageAsBase64(screenshots[currentImageIndex]);
				setLoadedImages(prev => new Map(prev).set(currentImageIndex, base64Image));
			} catch (error) {
				console.error(`Failed to load image ${currentImageIndex} for ${id}:`, error);
			} finally {
				setLoading(false);
				setLoadingImage(false);
			}
		};

		loadCurrentImage();
	}, [id, currentImageIndex, screenshots, loadedImages]);

	// Preload next image
	useEffect(() => {
		if (screenshots.length <= 1) return;

		const nextIndex = (currentImageIndex + 1) % screenshots.length;
		if (loadedImages.has(nextIndex)) return;

		const preloadNext = async () => {
			try {
				const base64Image = await GetImageAsBase64(screenshots[nextIndex]);
				setLoadedImages(prev => new Map(prev).set(nextIndex, base64Image));
			} catch (error) {
				console.error(`Failed to preload image ${nextIndex} for ${id}:`, error);
			}
		};

		preloadNext();
	}, [currentImageIndex, screenshots, loadedImages, id]);

	const handlePrevious = (e: React.MouseEvent) => {
		e.stopPropagation();
		setCurrentImageIndex((prev) =>
			prev === 0 ? screenshots.length - 1 : prev - 1
		);
	};

	const handleNext = (e: React.MouseEvent) => {
		e.stopPropagation();
		setCurrentImageIndex((prev) =>
			(prev + 1) % screenshots.length
		);
	};

	const currentImage = loadedImages.get(currentImageIndex);

	return (
		<Card
			className={`overflow-hidden hover:shadow-lg transition-shadow ${
				onClick ? "cursor-pointer" : ""
			}`}
			onClick={onClick}
		>
			<CardContent className="p-0">
				{/* Screenshot Carousel */}
				<div className="aspect-square bg-gray-100 relative group">
					{currentImage ? (
						<img
							src={currentImage}
							alt={`${name} - ${currentImageIndex + 1}`}
							className="w-full h-full object-cover"
						/>
					) : loading || loadingImage ? (
						<div className="w-full h-full flex items-center justify-center text-muted-foreground">
							<div className="flex flex-col items-center gap-2">
								<div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900" />
								<span className="text-xs">Loading...</span>
							</div>
						</div>
					) : screenshots.length > 0 ? (
						<div className="w-full h-full flex items-center justify-center text-muted-foreground">
							<div className="text-center">
								<div className="text-sm">Failed to load</div>
								<div className="text-xs">image</div>
							</div>
						</div>
					) : (
						<div className="w-full h-full flex items-center justify-center text-muted-foreground">
							No screenshot
						</div>
					)}

					{/* ID Badge */}
					<div className="absolute top-2 left-2 bg-black/70 text-white text-xs px-2 py-1 rounded">
						ID: {id}
					</div>

					{/* Image Counter */}
					{screenshots.length > 1 && (
						<div className="absolute top-2 right-2 bg-black/70 text-white text-xs px-2 py-1 rounded">
							{currentImageIndex + 1} / {screenshots.length}
						</div>
					)}

					{/* Carousel Navigation - Show on hover */}
					{screenshots.length > 1 && currentImage && (
						<>
							<Button
								variant="ghost"
								size="icon"
								className="absolute left-2 top-1/2 -translate-y-1/2 bg-black/50 hover:bg-black/70 text-white opacity-0 group-hover:opacity-100 transition-opacity"
								onClick={handlePrevious}
							>
								<ChevronLeft className="h-6 w-6" />
							</Button>
							<Button
								variant="ghost"
								size="icon"
								className="absolute right-2 top-1/2 -translate-y-1/2 bg-black/50 hover:bg-black/70 text-white opacity-0 group-hover:opacity-100 transition-opacity"
								onClick={handleNext}
							>
								<ChevronRight className="h-6 w-6" />
							</Button>

							{/* Dots Indicator */}
							<div className="absolute bottom-2 left-1/2 -translate-x-1/2 flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
								{screenshots.map((_, index) => (
									<button
										type="button"
										key={index}
										onClick={(e) => {
											e.stopPropagation();
											setCurrentImageIndex(index);
										}}
										className={`w-2 h-2 rounded-full transition-all ${
											index === currentImageIndex
												? "bg-white w-4"
												: "bg-white/50 hover:bg-white/75"
										}`}
										aria-label={`Go to image ${index + 1}`}
									/>
								))}
							</div>
						</>
					)}
				</div>

				{/* Content Info */}
				<div className="p-4">
					<h3 className="font-semibold truncate" title={name}>
						{name}
					</h3>
					<p className="text-xs text-muted-foreground mt-1 line-clamp-2">
						{description || "No description available"}
					</p>
					<div className="mt-2 flex items-center gap-2 text-xs text-muted-foreground">
						<span>{screenshots.length} screenshot{screenshots.length !== 1 ? 's' : ''}</span>
					</div>
				</div>
			</CardContent>
		</Card>
	);
}
