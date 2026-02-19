import React from "react";
import {
    View,
    Text,
    StyleSheet,
    FlatList,
    ActivityIndicator,
    RefreshControl,
} from "react-native";
import { normalize } from "@/core/normalize";
import { useBooks } from "@/core/hooks/useBooks";

interface BookListProps {
    filters?: {
        genre?: string;
        author?: string;
        publisher?: string;
        shelf_id?: string;
        search?: string;
    };
}

export const BookList = ({ filters }: BookListProps) => {
    const {
        data,
        fetchNextPage,
        hasNextPage,
        isFetchingNextPage,
        isLoading,
        refetch,
        isRefetching,
    } = useBooks(filters);

    const books = data?.pages.flatMap((page) => page.data) || [];

    const renderItem = ({ item }: { item: any }) => (
        <View style={styles.card}>
            <View style={styles.bookInfo}>
                <Text style={styles.title}>{item.title}</Text>
                <Text style={styles.author}>{item.author}</Text>
                <Text style={styles.publisher}>{item.publisher}</Text>

                <Text style={styles.genre}>Genre: {item.genre}</Text>
            </View>
            <View style={styles.locationContainer}>
                <Text style={styles.locationLabel}>Location</Text>
                <Text style={styles.location}>{item.shelf_address}</Text>
                <Text style={styles.location}>Row: {item.row}, Col: {item.column}</Text>
            </View>
        </View>
    );

    const renderFooter = () => {
        if (!isFetchingNextPage) return null;
        return (
            <View style={styles.footerLoader}>
                <ActivityIndicator size="small" color="#0000ff" />
            </View>
        );
    };

    if (isLoading) {
        return (
            <View style={styles.center}>
                <ActivityIndicator size="large" color="#0000ff" />
            </View>
        );
    }

    return (
        <View style={styles.container}>

            <FlatList
                data={books}
                renderItem={renderItem}
                keyExtractor={(item, index) => `${item._id}-${index}`}
                contentContainerStyle={styles.listContent}
                onEndReached={() => {
                    if (hasNextPage) fetchNextPage();
                }}
                onEndReachedThreshold={0.5}
                ListFooterComponent={renderFooter}
                refreshControl={
                    <RefreshControl refreshing={isRefetching} onRefresh={refetch} />
                }
                ListEmptyComponent={
                    <View style={styles.center}>
                        <Text style={styles.emptyText}>No books found.</Text>
                    </View>
                }
            />
        </View>
    );
};

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: "#F5F7FA",
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
        marginBottom: normalize(2),
    },

    genre: {
        fontSize: normalize(12),
        color: "#4A5568",
        marginTop: normalize(2),
        fontStyle: "italic",
    },
    footerLoader: {
        paddingVertical: normalize(20),
        alignItems: "center",
    },
    emptyText: {
        fontSize: normalize(16),
        color: "#718096",
    },
    locationContainer: {
        alignItems: "flex-end",
    },
    locationLabel: {
        fontSize: normalize(10),
        color: "#A0AEC0",
        marginBottom: normalize(2),
    },
    location: {
        fontSize: normalize(12),
        fontWeight: "500",
        color: "#4A5568",
    },
});
