import { GetHeroAssets } from "../../wailsjs/go/main/App"


export type HeroAssets = {
  id: number
  name: string
}

export async function getHeroAssets(): Promise<HeroAssets[]> {
  return GetHeroAssets();
}