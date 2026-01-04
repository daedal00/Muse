import React from "react";

interface StarRatingProps {
  value?: number | null;
  onChange?: (value: number) => void;
  label?: string;
  size?: number;
  className?: string;
}

const StarRating: React.FC<StarRatingProps> = ({
  value,
  onChange,
  label,
  size = 20,
  className = "",
}) => {
  const [hoveredValue, setHoveredValue] = React.useState<number | null>(null);
  const safeValue = Math.max(0, Math.min(5, value ?? 0));
  const displayValue = onChange && hoveredValue !== null ? hoveredValue : safeValue;
  const roundedValue = onChange ? displayValue : Math.round(displayValue);

  return (
    <div className={`flex items-center gap-1 ${className}`}>
      {Array.from({ length: 5 }, (_, index) => {
        const starValue = index + 1;
        const isFilled = starValue <= roundedValue;
        const icon = (
          <svg
            width={size}
            height={size}
            viewBox="0 0 24 24"
            fill="currentColor"
            className={
              isFilled
                ? "text-amber-500"
                : "text-slate-300 dark:text-slate-700"
            }
            aria-hidden="true"
          >
            <path d="M12 17.27l5.18 3.04-1.4-5.96 4.64-4.02-6.08-.52L12 4l-2.34 5.83-6.08.52 4.64 4.02-1.4 5.96L12 17.27z" />
          </svg>
        );

        if (!onChange) {
          return <span key={starValue}>{icon}</span>;
        }

        return (
          <button
            key={starValue}
            type="button"
            onClick={() => onChange(starValue)}
            onMouseEnter={() => setHoveredValue(starValue)}
            onMouseLeave={() => setHoveredValue(null)}
            aria-label={`Set rating to ${starValue} out of 5`}
            className="transition-transform hover:scale-105"
          >
            {icon}
          </button>
        );
      })}
      {label && <span className="ml-2 muted">{label}</span>}
    </div>
  );
};

export default StarRating;
