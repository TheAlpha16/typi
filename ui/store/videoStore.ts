import { create } from "zustand";
import { persist } from "zustand/middleware";

interface Video {
    vid: string;
    title: string;
    description: string;
    published_at: string;
    thumbnail: string;
}

interface VideoState {
    videos: { [page: number]: Video[] }; // cached videos
    currentPage: number; // current page number
    totalPages: number; // total number of pages
    isLoading: boolean; 
    searchQuery: string; // keyword to search videos
    filteredVideos: Video[]; // videos to be shown on the UI
    fetchVideos: (page: number, perPage: number) => Promise<void>; // fetch videos from the API
    setSearchQuery: (query: string) => void;
}

export const useVideoStore = create<VideoState>()(
    persist(
        (set, get) => ({
            videos: {},
            currentPage: 1,
            totalPages: 1,
            isLoading: true,
            searchQuery: "",
            filteredVideos: [],
            fetchVideos: async (page: number, perPage: number) => {
                set({ isLoading: true });
                try {
                    const res = await fetch(`/api/videos?page=${page}&per_page=${perPage}`);
                    const data = await res.json();
                    set((state) => {
                        const updatedVideos = { ...state.videos, [page]: data.videos };
                        return {
                            videos: updatedVideos,
                            currentPage: page,
                            totalPages: data.total_pages,
                            isLoading: false,
                            filteredVideos: get().searchQuery
                                ? filterVideos(updatedVideos, get().searchQuery)
                                : flattenVideos(updatedVideos),
                        };
                    });
                } catch (error) {
                    console.log("Failed to fetch videos:", error);
                    set({ isLoading: false });
                }
            },
            setSearchQuery: (query: string) => {
                set({
                    searchQuery: query,
                    filteredVideos: filterVideos(get().videos, query),
                });
            },
        }),
        {
            name: "video-storage",
        }
    )
);

// Helper function to filter videos based on query
const filterVideos = (videos: { [page: number]: Video[] }, query: string): Video[] => {
    const search = query.toLowerCase();
    return flattenVideos(videos).filter(
        (video) =>
            video.title.toLowerCase().includes(search) ||
            video.description.toLowerCase().includes(search)
    );
};

// Helper function to flatten paginated videos into a single array
const flattenVideos = (videos: { [page: number]: Video[] }): Video[] => {
    return Object.values(videos).flat();
};
