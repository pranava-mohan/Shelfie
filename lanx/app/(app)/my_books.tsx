import React, { useEffect, useState } from "react";
import {
    View,
    Text,
    StyleSheet,
    FlatList,
    ActivityIndicator,
    RefreshControl,
    Platform,
    StatusBar,
} from "react-native";
import { useClient } from "@/core/hooks/useClient";
import { API } from "@/config";
import { useAppSelector } from "@/core/store/store";
import { normalize } from "@/core/normalize";
import { SafeAreaView } from "react-native-safe-area-context";

interface Book {
    id: string;
    title: string;
    author: string;
    publisher: string;
    isbn: string;
}

interface BorrowedBookResult {
    book_details: Book;
    issued_at: string;
}

export default function MyBooksPage() {
    const { client } = useClient();
    const [books, setBooks] = useState<BorrowedBookResult[]>([]);
    const [loading, setLoading] = useState(true);
    const [refreshing, setRefreshing] = useState(false);
    const userID = useAppSelector((state) => state.auth.id);

    const fetchBooks = async () => {
        try {
            const response = await client.post(`${API.BASE_URL}/user/my-books`);
            setBooks(response.data);
        } catch (error) {
            console.error("Failed to fetch books:", error);
        } finally {
            setLoading(false);
            setRefreshing(false);
        }
    };

    useEffect(() => {
        fetchBooks();
    }, []);

    const onRefresh = () => {
        setRefreshing(true);
        fetchBooks();
    };

    const renderItem = ({ item }: { item: BorrowedBookResult }) => {
        const date = new Date(item.issued_at).toLocaleDateString();
        return (
            <View style={styles.card}>
                <View style={styles.bookInfo}>
                    <Text style={styles.title}>{item.book_details.title}</Text>
                    <Text style={styles.author}>{item.book_details.author}</Text>
                    <Text style={styles.publisher}>{item.book_details.publisher}</Text>
                </View>
                <View style={styles.dateContainer}>
                    <Text style={styles.dateLabel}>Issued on</Text>
                    <Text style={styles.date}>{date}</Text>
                </View>
            </View>
        );
    };

    if (loading) {
        return (
            <View style={styles.center}>
                <ActivityIndicator size="large" color="#0000ff" />
            </View>
        );
    }

    return (
        <SafeAreaView style={styles.safeArea}>
            <View style={styles.container}>
                <View style={styles.header}>
                    <Text style={styles.headerTitle}>My Books</Text>
                </View>
                <FlatList
                    data={books}
                    renderItem={renderItem}
                    keyExtractor={(item, index) => item.book_details.id || index.toString()}
                    contentContainerStyle={styles.listContent}
                    refreshControl={
                        <RefreshControl refreshing={refreshing} onRefresh={onRefresh} />
                    }
                    ListEmptyComponent={
                        <View style={styles.center}>
                            <Text style={styles.emptyText}>You haven't issued any books yet.</Text>
                        </View>
                    }
                />
            </View>
        </SafeAreaView>
    );
}

const styles = StyleSheet.create({
    safeArea: {
        flex: 1,
        backgroundColor: "#F5F7FA",
    },
    container: {
        flex: 1,
    },
    header: {
        paddingHorizontal: normalize(16),
        paddingVertical: normalize(12),
        backgroundColor: "#F5F7FA",
        borderBottomWidth: 1,
        borderBottomColor: "#E2E8F0",
    },
    headerTitle: {
        fontSize: normalize(24),
        fontWeight: "bold",
        color: "#1A202C",
    },
    center: {
        flex: 1,
        justifyContent: "center",
        alignItems: "center",
        paddingTop: normalize(50),
    },
    listContent: {
        padding: normalize(16),
        paddingBottom: normalize(100),
    },
    card: {
        backgroundColor: "#FFFFFF",
        borderRadius: normalize(12),
        padding: normalize(16),
        marginBottom: normalize(12),
        shadowColor: "#000",
        shadowOffset: {
            width: 0,
            height: 2,
        },
        shadowOpacity: 0.1,
        shadowRadius: 4,
        elevation: 3,
        flexDirection: "row",
        justifyContent: "space-between",
        alignItems: "center",
    },
    bookInfo: {
        flex: 1,
        marginRight: normalize(10),
    },
    title: {
        fontSize: normalize(16),
        fontWeight: "600",
        color: "#2D3748",
        marginBottom: normalize(4),
    },
    author: {
        fontSize: normalize(14),
        color: "#718096",
        marginBottom: normalize(2),
    },
    publisher: {
        fontSize: normalize(12),
        color: "#A0AEC0",
    },
    dateContainer: {
        alignItems: "flex-end",
        marginBottom: normalize(14)
    },
    dateLabel: {
        fontSize: normalize(10),
        color: "#A0AEC0",
        marginBottom: normalize(2),
    },
    date: {
        fontSize: normalize(12),
        fontWeight: "500",
        color: "#4A5568",
    },
    emptyText: {
        fontSize: normalize(16),
        color: "#718096",
    },
});
