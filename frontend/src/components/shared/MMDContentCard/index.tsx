import { useState, useEffect } from "react";
import { GetImageAsBase64, GetVideoAsBase64 } from "../../../../wailsjs/go/handlers/Images";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { ChevronLeft, ChevronRight } from "lucide-react";

interface MMDContentCardProps {
	id: string;
	name: string;
	screenshots: string[] | null;
	video?: string[] | null;
	description: string;
	onClick?: () => void;
}

export function MMDContentCard({
	id,
	name,
	screenshots,
	video,
	description,
	onClick,
}: MMDContentCardProps) {
	const [currentImageIndex, setCurrentImageIndex] = useState(0);
	const [loadedImages, setLoadedImages] = useState<Map<number, string>>(new Map());
	const [loading, setLoading] = useState(true);
	const [loadingImage, setLoadingImage] = useState(false);
	const [videoSrc, setVideoSrc] = useState<string>("");

	// Normalize null to empty array for easier handling
	const normalizedScreenshots = screenshots ?? [];
	const normalizedVideo = video ?? [];

	// Determine if we should show video (prefer video over screenshots)
	const hasVideo = normalizedVideo.length > 0;
	const hasScreenshots = normalizedScreenshots.length > 0;

	// Load video (if available)
	useEffect(() => {
		if (!hasVideo) return;

		const loadVideo = async () => {
			setLoading(true);
			try {
				const base64Video = await GetVideoAsBase64(normalizedVideo[0]);
				setVideoSrc(base64Video);
			} catch (error) {
				console.error(`Failed to load video for ${id}:`, error);
			} finally {
				setLoading(false);
			}
		};

		loadVideo();
	}, [hasVideo, normalizedVideo, id]);

	// Load the current image (only if no video)
	useEffect(() => {
		if (hasVideo) {
			setLoading(false);
			return;
		}

		const loadCurrentImage = async () => {
			if (normalizedScreenshots.length === 0) {
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
				const base64Image = await GetImageAsBase64(normalizedScreenshots[currentImageIndex]);
				setLoadedImages(prev => new Map(prev).set(currentImageIndex, base64Image));
			} catch (error) {
				console.error(`Failed to load image ${currentImageIndex} for ${id}:`, error);
			} finally {
				setLoading(false);
				setLoadingImage(false);
			}
		};

		loadCurrentImage();
	}, [id, currentImageIndex, normalizedScreenshots, loadedImages, hasVideo]);

	// Preload next image
	useEffect(() => {
		if (normalizedScreenshots.length <= 1) return;

		const nextIndex = (currentImageIndex + 1) % normalizedScreenshots.length;
		if (loadedImages.has(nextIndex)) return;

		const preloadNext = async () => {
			try {
				const base64Image = await GetImageAsBase64(normalizedScreenshots[nextIndex]);
				setLoadedImages(prev => new Map(prev).set(nextIndex, base64Image));
			} catch (error) {
				console.error(`Failed to preload image ${nextIndex} for ${id}:`, error);
			}
		};

		preloadNext();
	}, [currentImageIndex, normalizedScreenshots, loadedImages, id]);

	const handlePrevious = (e: React.MouseEvent) => {
		e.stopPropagation();
		setCurrentImageIndex((prev) =>
			prev === 0 ? normalizedScreenshots.length - 1 : prev - 1
		);
	};

	const handleNext = (e: React.MouseEvent) => {
		e.stopPropagation();
		setCurrentImageIndex((prev) =>
			(prev + 1) % normalizedScreenshots.length
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
				{/* Media Display - Video or Screenshot Carousel */}
				<div className="aspect-square bg-gray-100 relative group">
					{hasVideo ? (
						/* Show first video when available */
						videoSrc ? (
							<video
								src={videoSrc}
								controls
								className="w-full h-full object-cover"
								preload="metadata"
							>
								Your browser does not support the video tag.
							</video>
						) : (
							<div className="w-full h-full flex items-center justify-center text-muted-foreground">
								<div className="flex flex-col items-center gap-2">
									<div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900" />
									<span className="text-xs">Loading video...</span>
								</div>
							</div>
						)
					) : currentImage ? (
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
					) : hasScreenshots ? (
						<div className="w-full h-full flex items-center justify-center text-muted-foreground">
							<div className="text-center">
								<div className="text-sm">Failed to load</div>
								<div className="text-xs">image</div>
							</div>
						</div>
					) : (
						<div className="w-full h-full flex items-center justify-center text-muted-foreground">
							<div className="text-center text-sm">No media available</div>
						</div>
					)}

					{/* ID Badge */}
					<div className="absolute top-2 left-2 bg-black/70 text-white text-xs px-2 py-1 rounded">
						ID: {id}
					</div>

					{/* Media Counter */}
					{hasVideo && normalizedVideo.length > 1 && (
						<div className="absolute top-2 right-2 bg-black/70 text-white text-xs px-2 py-1 rounded">
							{normalizedVideo.length} video{normalizedVideo.length !== 1 ? 's' : ''}
						</div>
					)}
					{!hasVideo && normalizedScreenshots.length > 1 && (
						<div className="absolute top-2 right-2 bg-black/70 text-white text-xs px-2 py-1 rounded">
							{currentImageIndex + 1} / {normalizedScreenshots.length}
						</div>
					)}

					{/* Carousel Navigation - Show on hover (only for screenshots) */}
					{!hasVideo && normalizedScreenshots.length > 1 && currentImage && (
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
								{normalizedScreenshots.map((_, index) => (
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
						{hasVideo && (
							<span>{normalizedVideo.length} video{normalizedVideo.length !== 1 ? 's' : ''}</span>
						)}
						{hasScreenshots && (
							<span>{normalizedScreenshots.length} screenshot{normalizedScreenshots.length !== 1 ? 's' : ''}</span>
						)}
						{!hasVideo && !hasScreenshots && (
							<span>No media</span>
						)}
					</div>
				</div>
			</CardContent>
		</Card>
	);
}
