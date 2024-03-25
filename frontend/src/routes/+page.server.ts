import { listWines } from '$lib/api/wines';
import { error } from '@sveltejs/kit';
import type { ServerLoadEvent } from '@sveltejs/kit';


/** @type {import('./$types').PageServerLoad} */
export async function load(_: ServerLoadEvent) {
	const wines = await listWines();
	const images: string[] = [
		"https://images.vivino.com/thumbs/uDMfPG10R1efm5H1do--Ow_pb_600x600.png",
		"https://images.vivino.com/thumbs/Eh506eKdSXGvndyHHfqPug_pb_600x600.png",
		"https://images.vivino.com/thumbs/MCs7Ix2zS56U3-vs6GOR4A_pb_600x600.png",
	];

	if (wines) {
		return { wines, images };
	}

	error(404, 'Not found');
}

export const ssr = true;
