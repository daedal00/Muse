import React, { useState } from "react";
import { useMutation } from "@apollo/client";
import { CREATE_COMMENT } from "../lib/graphql/mutations";
import Link from "next/link";

interface Comment {
  id: string;
  content: string;
  createdAt: string;
  user: {
    id: string;
    name: string;
    avatar?: string;
  };
}

interface CommentSectionProps {
  targetId: string;
  targetType: "ALBUM" | "TRACK";
  comments: { edges: { node: Comment }[]; totalCount: number };
  isAuthed: boolean;
  refetchQueries?: any[];
}

const CommentSection: React.FC<CommentSectionProps> = ({
  targetId,
  targetType,
  comments,
  isAuthed,
  refetchQueries,
}) => {
  const [content, setContent] = useState("");
  const [createComment, { loading, error }] = useMutation(CREATE_COMMENT, {
    refetchQueries,
    awaitRefetchQueries: true,
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!content.trim()) return;

    try {
      await createComment({
        variables: {
          input: {
            content: content.trim(),
            albumId: targetType === "ALBUM" ? targetId : undefined,
            trackId: targetType === "TRACK" ? targetId : undefined,
          },
        },
      });
      setContent("");
    } catch (e) {
      console.error(e);
    }
  };

  return (
    <div className="card">
      <h2 className="text-2xl font-semibold mb-4">Comments</h2>

      {/* List */}
      <div className="space-y-6 mb-8">
        {comments && comments.edges.length > 0 ? (
          comments.edges.map(({ node }) => (
            <div key={node.id} className="flex gap-4">
              <div className="flex-shrink-0">
                {node.user.avatar ? (
                  <img
                    src={node.user.avatar}
                    className="w-10 h-10 rounded-full object-cover"
                  />
                ) : (
                  <div className="w-10 h-10 rounded-full bg-slate-200 dark:bg-slate-700 flex items-center justify-center text-slate-500 font-bold">
                    {node.user.name[0]?.toUpperCase()}
                  </div>
                )}
              </div>
              <div className="flex-1">
                <div className="flex items-baseline justify-between">
                  <span className="font-semibold text-slate-900 dark:text-slate-100">
                    {node.user.name}
                  </span>
                  <span className="text-xs text-slate-500">
                    {new Date(node.createdAt).toLocaleDateString()}{" "}
                    {new Date(node.createdAt).toLocaleTimeString([], {
                      hour: "2-digit",
                      minute: "2-digit",
                    })}
                  </span>
                </div>
                <p className="text-slate-600 dark:text-slate-300 mt-1 text-sm whitespace-pre-wrap">
                  {node.content}
                </p>
              </div>
            </div>
          ))
        ) : (
          <p className="text-slate-500 italic">
            No comments yet. Be the first!
          </p>
        )}
      </div>

      {/* Form */}
      {isAuthed ? (
        <form onSubmit={handleSubmit} className="relative">
          <textarea
            value={content}
            onChange={(e) => setContent(e.target.value)}
            placeholder="Join the discussion..."
            className="w-full p-4 rounded-xl bg-slate-50 dark:bg-slate-800/50 border border-slate-200 dark:border-slate-700 focus:outline-none focus:ring-2 focus:ring-emerald-500/50 resize-y min-h-[100px] text-slate-800 dark:text-slate-200"
          />
          <div className="mt-2 flex justify-end">
            <button
              type="submit"
              disabled={loading || !content.trim()}
              className="btn-primary"
            >
              {loading ? "Posting..." : "Post Comment"}
            </button>
          </div>
          {error && (
            <p className="text-rose-500 text-sm mt-2">{error.message}</p>
          )}
        </form>
      ) : (
        <div className="p-4 rounded-xl bg-slate-50 dark:bg-slate-800/50 text-center">
          <p className="text-slate-500 mb-2">Log in to join the conversation</p>
          <Link
            href="/auth"
            className="text-emerald-600 hover:text-emerald-500 font-medium"
          >
            Login
          </Link>
        </div>
      )}
    </div>
  );
};

export default CommentSection;
