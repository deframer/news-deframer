export type HostStatus = 'idle' | 'loading' | 'success' | 'error';

interface StatusBadgeProps {
  status: HostStatus;
  enabled?: boolean;
  labels: {
    connected: string;
    error: string;
    checking: string;
    disabled: string;
  };
}

export const StatusBadge = ({ status, enabled = true, labels }: StatusBadgeProps) => {
  if (!enabled) {
    return (
      <div className="status-badge" role="status" aria-live="polite">
        <span className="status-dot" />
        {labels.disabled}
      </div>
    );
  }

  const className = status === 'success' ? 'connected' : status === 'error' ? 'error' : status === 'loading' ? 'testing' : '';
  const label = status === 'success' ? labels.connected : status === 'error' ? labels.error : labels.checking;

  return (
    <div className={`status-badge ${className}`.trim()} role="status" aria-live="polite">
      <span className="status-dot" />
      {label}
    </div>
  );
};
