import { useInfiniteQuery } from "@tanstack/react-query";
import { useClient } from "./useClient";

interface Book {
    _id: string; // Assuming ObjectId is hex string
    title: string;
    author: string;
    publisher: string;
    isbn: string;
    shelf_id: string;
    row: number;
    column: number;
    shelf_address: string;
    added_at: string;
    genre: string;
}

interface FetchBooksResponse {
    data: Book[];
    meta: {
        total: number;
        page: number;
        last_page: number;
        limit: number;
    };
}

interface UseBooksParams {
    genre?: string;
    author?: string;
    publisher?: string;
    shelf_id?: string;
    search?: string;
}

export const useBooks = (filters: UseBooksParams = {}) => {
    const { client } = useClient();

    const fetchBooks = async ({ pageParam = 1 }) => {
        const { data } = await client.post<FetchBooksResponse>("/book/all", {
            page: pageParam,
            limit: 10,
            available_only: true,
            ...filters
        });
        return data;
    };

    return useInfiniteQuery({
        queryKey: ["books", filters],
        queryFn: fetchBooks,
        getNextPageParam: (lastPage) => {
            if (lastPage.meta.page < lastPage.meta.last_page) {
                return lastPage.meta.page + 1;
            }
            return undefined;
        },
        initialPageParam: 1,
    });
};
