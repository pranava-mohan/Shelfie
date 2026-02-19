import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';

export const useAdminAuth = () => {
    const router = useRouter();
    const [token, setToken] = useState<string | null>(null);

    useEffect(() => {
        const storedToken = localStorage.getItem('admin_token');
        if (!storedToken) {
            router.replace('/admin/login');
        } else {
            setToken(storedToken);
        }
    }, [router]);

    const logout = () => {
        localStorage.removeItem('admin_token');
        router.replace('/admin/login');
    };

    return { token, logout };
};
