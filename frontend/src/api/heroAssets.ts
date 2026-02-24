export type HeroAssets = {
  id: number
  name: string
}

export async function getHeroAssets(): Promise<HeroAssets[]> {
  const res = await fetch('http://localhost:3000/get-hero-assets');
  if (!res.ok) {
    throw new Error('Failed to fetch hero assets')
  }
  const data = await res.json()
  console.log(data)
  return data
}
