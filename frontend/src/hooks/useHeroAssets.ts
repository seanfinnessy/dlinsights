import { useQuery } from "@tanstack/react-query";
import { getHeroAssets } from "../api/heroAssets";

export function useHeroAssets() {
    return useQuery({
        queryKey: ['heroAssets'],
        queryFn: getHeroAssets,
    })
}