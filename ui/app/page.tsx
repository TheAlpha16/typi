"use client";

import React, { useEffect, useRef, useCallback } from "react";
import { useInView } from "react-intersection-observer";
import { useVideoStore } from "@/store/videoStore";
import VideoCard from "@/components/VideoCard";
import VideoSkeleton from "@/components/VideoSkeleton";
import { Youtube, Search } from "lucide-react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { ThemeToggle } from "@/components/ThemeProvider";

export default function Home() {
  const {
    filteredVideos,
    currentPage,
    totalPages,
    isLoading,
    searchQuery,
    fetchVideos,
    setSearchQuery,
  } = useVideoStore();
  const loadedPages = useRef(new Set<number>());

  const { ref, inView } = useInView({ threshold: 0 });

  const loadMoreVideos = useCallback(() => {
    if (
      currentPage < totalPages &&
      !isLoading &&
      !loadedPages.current.has(currentPage + 1)
    ) {
      fetchVideos(currentPage + 1, 10);
      loadedPages.current.add(currentPage + 1);
    }
  }, [currentPage, totalPages, isLoading, fetchVideos]);

  useEffect(() => {
    if (inView) {
      loadMoreVideos();
    }
  }, [inView, loadMoreVideos]);

  useEffect(() => {
    if (currentPage === 1 && filteredVideos.length === 0) {
      fetchVideos(1, 10);
      loadedPages.current.add(1);
    }
  }, [fetchVideos, filteredVideos, currentPage]);

  const genRandomKey = () => Math.random().toString(36).substring(7);

  return (
    <div className="min-h-screen bg-gray-100 dark:bg-gray-900 transition-colors duration-200">
      <div className="container mx-auto px-4 py-8">
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-3xl font-bold text-gray-800 dark:text-gray-100 flex items-center">
            <Youtube className="mr-2" /> typi
          </h1>
          <ThemeToggle />
        </div>
        <span className="flex text-gray-600 dark:text-gray-400 mb-6">
          Typi fetches videos from the API in paginated responses, enabling lazy loading and infinite scrolling for a seamless browsing experience.
        </span>
        <div className="mb-6 flex flex-col sm:flex-row gap-4">
          <form
            onSubmit={(e) => e.preventDefault()}
            className="flex-grow flex gap-2"
          >
            <Input
              type="text"
              placeholder="filter cricket videos with keywords. eg: rohit sharma"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="flex-grow truncate"
            />
            <Button type="submit">
              <Search size={"icon"} />
            </Button>
          </form>
        </div>
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
          {filteredVideos.map((video) => (
            <VideoCard key={genRandomKey()} video={video} />
          ))}
          {isLoading &&
            Array.from({ length: 10 }).map((_, index) => (
              <VideoSkeleton key={index} />
            ))}
        </div>
        <div ref={ref} className="h-10" />
      </div>
    </div>
  );
}
