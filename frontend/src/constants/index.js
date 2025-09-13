// Use a single base URL for all API calls
const API_BASE_URL = process.env.REACT_APP_API_BASE_URL || "https://hotel.i-dmytro.org/api";

export const PRICE_API_URL = `${API_BASE_URL}/prices/`;
export const PROGRAM_API_URL = `${API_BASE_URL}/programs/`;
export const NEWS_API_URL = `${API_BASE_URL}/news/`;