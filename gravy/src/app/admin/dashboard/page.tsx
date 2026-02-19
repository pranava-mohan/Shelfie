'use client';
import { useEffect, useState } from 'react';
import axios from 'axios';
import { useAdminAuth } from '@/hooks/useAdminAuth';
import { API_URL } from '@/config';
import toast, { Toaster } from 'react-hot-toast';
import { PlusIcon, PencilIcon, TrashIcon, ArchiveBoxIcon, ArrowDownTrayIcon } from '@heroicons/react/24/solid';
import QRCode from 'qrcode';

type Book = {
    id: string;
    title: string;
    author: string;
    publisher: string;
    isbn: string;
    shelf_id: string;
    row: number;
    column: number;
};

type Shelf = {
    id: string;
    address: string;
};

export default function AdminDashboard() {
    const { token, logout } = useAdminAuth();
    const [books, setBooks] = useState<Book[]>([]);
    const [shelves, setShelves] = useState<Shelf[]>([]);
    const [isBookModalOpen, setIsBookModalOpen] = useState(false);
    const [isShelfModalOpen, setIsShelfModalOpen] = useState(false);
    const [isShelfListOpen, setIsShelfListOpen] = useState(false);
    const [isEditing, setIsEditing] = useState(false);
    const [currentBook, setCurrentBook] = useState<Partial<Book>>({});
    const [newShelfAddress, setNewShelfAddress] = useState('');
    const [currentShelfId, setCurrentShelfId] = useState<string | null>(null);
    const [page, setPage] = useState(1);
    const [meta, setMeta] = useState<{ total: number, last_page: number } | null>(null);

    useEffect(() => {
        if (token) {
            fetchBooks(1);
            fetchShelves();
        }
    }, [token]);

    const fetchBooks = async (pageNum = 1) => {
        try {
            const { data } = await axios.post(
                `${API_URL}/book/all`,
                { page: pageNum, limit: 10 },
                { headers: { Authorization: `Bearer ${token}` } }
            );
            if (data.data) {
                setBooks(data.data);
                setMeta(data.meta);
                setPage(data.meta.page);
            } else {
                setBooks(data);
            }
        } catch (error) {
            console.error('Failed to fetch books', error);
            toast.error('Failed to load books');
        }
    };

    const fetchShelves = async () => {
        try {
            const { data } = await axios.post(
                `${API_URL}/shelf/all`,
                {},
                { headers: { Authorization: `Bearer ${token}` } }
            );
            setShelves(data);
        } catch (error) {
            console.error('Failed to fetch shelves', error);
        }
    };

    const handleDeleteBook = async (id: string) => {
        if (!confirm('Are you sure you want to delete this book?')) return;
        try {
            await axios.post(
                `${API_URL}/book/delete`,
                { book_id: id },
                { headers: { Authorization: `Bearer ${token}` } }
            );
            toast.success('Book deleted');
            fetchBooks(page);
        } catch (error) {
            toast.error('Failed to delete book');
        }
    };

    const handleSaveBook = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            if (isEditing && currentBook.id) {
                await axios.post(
                    `${API_URL}/book/update`,
                    {
                        book_id: currentBook.id,
                        ...currentBook
                    },
                    { headers: { Authorization: `Bearer ${token}` } }
                );
                toast.success('Book updated');
            } else {
                await axios.post(
                    `${API_URL}/book/create`,
                    currentBook,
                    { headers: { Authorization: `Bearer ${token}` } }
                );
                toast.success('Book created');
            }
            setIsBookModalOpen(false);
            fetchBooks(page);
        } catch (error: any) {
            toast.error(error.response?.data?.error || 'Failed to save book');
        }
    };

    const handleSaveShelf = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            if (currentShelfId) {
                await axios.post(
                    `${API_URL}/shelf/update`,
                    { shelf_id: currentShelfId, address: newShelfAddress },
                    { headers: { Authorization: `Bearer ${token}` } }
                );
                toast.success('Shelf updated');
            } else {
                await axios.post(
                    `${API_URL}/shelf/create`,
                    { address: newShelfAddress },
                    { headers: { Authorization: `Bearer ${token}` } }
                );
                toast.success('Shelf created');
            }
            setNewShelfAddress('');
            setCurrentShelfId(null);
            setIsShelfModalOpen(false);
            fetchShelves();
        } catch (error: any) {
            toast.error('Failed to save shelf');
        }
    };

    const handleDeleteShelf = async (id: string) => {
        if (!confirm('Are you sure? This might act weird if books are in it.')) return;
        try {
            await axios.post(
                `${API_URL}/shelf/delete`,
                { shelf_id: id },
                { headers: { Authorization: `Bearer ${token}` } }
            );
            toast.success('Shelf deleted');
            fetchShelves();
        } catch (error: any) {
            toast.error('Failed to delete shelf');
        }
    }

    const handleDownloadQR = async (book: Book) => {
        try {
            const url = await QRCode.toDataURL(book.id);
            const a = document.createElement('a');
            a.href = url;
            a.download = `book-${book.id}-qr.png`;
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);
            toast.success('QR Code downloaded');
        } catch (err) {
            console.error(err);
            toast.error('Failed to generate QR Code');
        }
    };


    const openAddBookModal = () => {
        setIsEditing(false);
        setCurrentBook({
            title: '',
            author: '',
            publisher: '',
            isbn: '',
            shelf_id: shelves.length > 0 ? shelves[0].id : '',
            row: 1,
            column: 1
        });
        setIsBookModalOpen(true);
    };

    const openEditBookModal = (book: Book) => {
        setIsEditing(true);
        setCurrentBook(book);
        setIsBookModalOpen(true);
    };

    const openAddShelfModal = () => {
        setNewShelfAddress('');
        setCurrentShelfId(null);
        setIsShelfModalOpen(true);
    }

    const openEditShelfModal = (shelf: Shelf) => {
        setNewShelfAddress(shelf.address);
        setCurrentShelfId(shelf.id);
        setIsShelfModalOpen(true);
    }

    return (
        <div className="min-h-screen bg-gray-50 p-6">
            <Toaster />
            <div className="mx-auto max-w-7xl">
                <div className="mb-6 flex flex-col items-start justify-between gap-4 sm:flex-row sm:items-center">
                    <h1 className="text-3xl font-bold text-gray-800">Library Admin</h1>
                    <div className="flex gap-4">
                        <button
                            onClick={() => setIsShelfListOpen(true)}
                            className="flex items-center gap-2 rounded-lg bg-gray-600 px-4 py-2 text-white hover:bg-gray-700"
                        >
                            <ArchiveBoxIcon className="h-5 w-5" /> Manage Shelves
                        </button>
                        <button
                            onClick={openAddBookModal}
                            className="flex items-center gap-2 rounded-lg bg-green-600 px-4 py-2 text-white hover:bg-green-700"
                        >
                            <PlusIcon className="h-5 w-5" /> Add Book
                        </button>
                        <button
                            onClick={logout}
                            className="rounded-lg bg-red-500 px-4 py-2 text-white hover:bg-red-600"
                        >
                            Logout
                        </button>
                    </div>
                </div>

                <div className="overflow-hidden rounded-lg bg-white shadow">
                    <table className="min-w-full divide-y divide-gray-200">
                        <thead className="bg-gray-50">
                            <tr>
                                <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">Title</th>
                                <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">Author</th>
                                <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">Shelf</th>
                                <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">Pos</th>
                                <th className="px-6 py-3 text-right text-xs font-medium uppercase tracking-wider text-gray-500">Actions</th>
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-gray-200 bg-white">
                            {books.map((book) => {
                                const shelf = shelves.find(s => s.id === book.shelf_id);
                                return (
                                    <tr key={book.id}>
                                        <td className="whitespace-nowrap px-6 py-4">{book.title}</td>
                                        <td className="whitespace-nowrap px-6 py-4">{book.author}</td>
                                        <td className="whitespace-nowrap px-6 py-4">
                                            {shelf ? shelf.address : <span className="text-gray-400 text-xs">{book.shelf_id}</span>}
                                        </td>
                                        <td className="whitespace-nowrap px-6 py-4">R{book.row}:C{book.column}</td>
                                        <td className="whitespace-nowrap px-6 py-4 text-right">
                                            <button onClick={() => openEditBookModal(book)} className="mr-2 text-blue-600 hover:text-blue-900">
                                                <PencilIcon className="h-5 w-5" />
                                            </button>
                                            <button onClick={() => handleDeleteBook(book.id)} className="mr-2 text-red-600 hover:text-red-900">
                                                <TrashIcon className="h-5 w-5" />
                                            </button>
                                            <button onClick={() => handleDownloadQR(book)} className="text-gray-600 hover:text-gray-900" title="Download QR Code">
                                                <ArrowDownTrayIcon className="h-5 w-5" />
                                            </button>
                                        </td>
                                    </tr>
                                )
                            })}
                            {books.length === 0 && (
                                <tr>
                                    <td colSpan={5} className="px-6 py-4 text-center text-gray-500">No books found.</td>
                                </tr>
                            )}
                        </tbody>
                    </table>
                    <div className="flex items-center justify-between border-t border-gray-200 bg-white px-4 py-3 sm:px-6">
                        <div className="flex flex-1 justify-between sm:hidden">
                            <button
                                onClick={() => fetchBooks(page - 1)}
                                disabled={page <= 1}
                                className="relative inline-flex items-center rounded-md border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50"
                            >
                                Previous
                            </button>
                            <button
                                onClick={() => fetchBooks(page + 1)}
                                disabled={!meta || page >= meta.last_page}
                                className="relative ml-3 inline-flex items-center rounded-md border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50"
                            >
                                Next
                            </button>
                        </div>
                        <div className="hidden sm:flex sm:flex-1 sm:items-center sm:justify-between">
                            <div>
                                <p className="text-sm text-gray-700">
                                    Showing page <span className="font-medium">{page}</span> of{' '}
                                    <span className="font-medium">{meta?.last_page || 1}</span> ({meta?.total || 0} results)
                                </p>
                            </div>
                            <div>
                                <nav className="isolate inline-flex -space-x-px rounded-md shadow-sm" aria-label="Pagination">
                                    <button
                                        onClick={() => fetchBooks(page - 1)}
                                        disabled={page <= 1}
                                        className="relative inline-flex items-center rounded-l-md px-2 py-2 text-gray-400 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-20 focus:outline-offset-0 disabled:opacity-50"
                                    >
                                        <span className="sr-only">Previous</span>
                                        Previous
                                    </button>
                                    <button
                                        onClick={() => fetchBooks(page + 1)}
                                        disabled={!meta || page >= meta.last_page}
                                        className="relative inline-flex items-center rounded-r-md px-2 py-2 text-gray-400 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-20 focus:outline-offset-0 disabled:opacity-50"
                                    >
                                        <span className="sr-only">Next</span>
                                        Next
                                    </button>
                                </nav>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            {isBookModalOpen && (
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
                    <div className="w-full max-w-lg rounded-lg bg-white p-6 shadow-xl max-h-[90vh] overflow-y-auto">
                        <h2 className="mb-4 text-xl font-bold">{isEditing ? 'Edit Book' : 'Add New Book'}</h2>
                        <form onSubmit={handleSaveBook} className="grid grid-cols-2 gap-4">
                            <div className="col-span-2">
                                <label className="block text-sm font-medium">Title</label>
                                <input
                                    className="mt-1 block w-full rounded border p-2"
                                    value={currentBook.title || ''}
                                    onChange={(e) => setCurrentBook({ ...currentBook, title: e.target.value })}
                                    required
                                />
                            </div>
                            <div className="col-span-2">
                                <label className="block text-sm font-medium">Author</label>
                                <input
                                    className="mt-1 block w-full rounded border p-2"
                                    value={currentBook.author || ''}
                                    onChange={(e) => setCurrentBook({ ...currentBook, author: e.target.value })}
                                    required
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-medium">Publisher</label>
                                <input
                                    className="mt-1 block w-full rounded border p-2"
                                    value={currentBook.publisher || ''}
                                    onChange={(e) => setCurrentBook({ ...currentBook, publisher: e.target.value })}
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-medium">ISBN</label>
                                <input
                                    className="mt-1 block w-full rounded border p-2"
                                    value={currentBook.isbn || ''}
                                    onChange={(e) => setCurrentBook({ ...currentBook, isbn: e.target.value })}
                                />
                            </div>
                            <div className="col-span-2">
                                <label className="block text-sm font-medium">Shelf</label>
                                <select
                                    className="mt-1 block w-full rounded border p-2 bg-white"
                                    value={currentBook.shelf_id || ''}
                                    onChange={(e) => setCurrentBook({ ...currentBook, shelf_id: e.target.value })}
                                    required
                                >
                                    <option value="" disabled>Select Shelf</option>
                                    {shelves.map(shelf => (
                                        <option key={shelf.id} value={shelf.id}>{shelf.address}</option>
                                    ))}
                                </select>
                            </div>
                            <div>
                                <label className="block text-sm font-medium">Row</label>
                                <input
                                    type="number"
                                    className="mt-1 block w-full rounded border p-2"
                                    value={currentBook.row || ''}
                                    onChange={(e) => setCurrentBook({ ...currentBook, row: parseInt(e.target.value) })}
                                    required
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-medium">Column</label>
                                <input
                                    type="number"
                                    className="mt-1 block w-full rounded border p-2"
                                    value={currentBook.column || ''}
                                    onChange={(e) => setCurrentBook({ ...currentBook, column: parseInt(e.target.value) })}
                                    required
                                />
                            </div>

                            <div className="col-span-2 flex justify-end gap-2 mt-4">
                                <button
                                    type="button"
                                    onClick={() => setIsBookModalOpen(false)}
                                    className="rounded bg-gray-300 px-4 py-2 hover:bg-gray-400"
                                >
                                    Cancel
                                </button>
                                <button
                                    type="submit"
                                    className="rounded bg-blue-600 px-4 py-2 text-white hover:bg-blue-700"
                                >
                                    Save
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            )}

            {isShelfListOpen && (
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
                    <div className="w-full max-w-lg rounded-lg bg-white p-6 shadow-xl">
                        <div className="flex justify-between items-center mb-4">
                            <h2 className="text-xl font-bold">Manage Shelves</h2>
                            <button onClick={openAddShelfModal} className="text-sm bg-green-100 text-green-800 px-2 py-1 rounded hover:bg-green-200">+ New Shelf</button>
                        </div>

                        <div className="max-h-[60vh] overflow-y-auto">
                            <table className="min-w-full divide-y divide-gray-200">
                                <thead className="bg-gray-50">
                                    <tr>
                                        <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Address</th>
                                        <th className="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase">Actions</th>
                                    </tr>
                                </thead>
                                <tbody className="divide-y divide-gray-200">
                                    {shelves.map(shelf => (
                                        <tr key={shelf.id}>
                                            <td className="px-4 py-3">{shelf.address}</td>
                                            <td className="px-4 py-3 text-right">
                                                <button onClick={() => openEditShelfModal(shelf)} className="mr-2 text-blue-500 hover:text-blue-700">Edit</button>
                                                <button onClick={() => handleDeleteShelf(shelf.id)} className="text-red-500 hover:text-red-700">Delete</button>
                                            </td>
                                        </tr>
                                    ))}
                                    {shelves.length === 0 && (
                                        <tr><td colSpan={2} className="px-4 py-2 text-center text-gray-400">No shelves found</td></tr>
                                    )}
                                </tbody>
                            </table>
                        </div>

                        <div className="mt-4 flex justify-end">
                            <button
                                onClick={() => setIsShelfListOpen(false)}
                                className="rounded bg-gray-300 px-4 py-2 hover:bg-gray-400"
                            >
                                Close
                            </button>
                        </div>
                    </div>
                </div>
            )}

            {isShelfModalOpen && (
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-60">
                    <div className="w-full max-w-sm rounded-lg bg-white p-6 shadow-xl">
                        <h2 className="mb-4 text-xl font-bold">{currentShelfId ? 'Edit Shelf' : 'Add New Shelf'}</h2>
                        <form onSubmit={handleSaveShelf}>
                            <div className="mb-4">
                                <label className="block text-sm font-medium">Shelf Address/Name</label>
                                <input
                                    className="mt-1 block w-full rounded border p-2"
                                    value={newShelfAddress}
                                    onChange={(e) => setNewShelfAddress(e.target.value)}
                                    placeholder="e.g. A1, Fiction Section"
                                    required
                                />
                            </div>
                            <div className="flex justify-end gap-2">
                                <button
                                    type="button"
                                    onClick={() => setIsShelfModalOpen(false)}
                                    className="rounded bg-gray-300 px-4 py-2 hover:bg-gray-400"
                                >
                                    Cancel
                                </button>
                                <button
                                    type="submit"
                                    className="rounded bg-green-600 px-4 py-2 text-white hover:bg-green-700"
                                >
                                    {currentShelfId ? 'Update' : 'Create'}
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            )}

        </div>
    );
}
