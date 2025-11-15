
type Props = {
    status: string;
    onChange: (newStatus: string) => void;
}

const statusMap = {
    "uncollected": "未回収",
    "collected": "回収済み",
    "lost": "消失"
}

const StatusPicker = (props: Props) => {

    return (
        <select value={props.status} className="border border-gray-300 p-2" onChange={(e) => {
            props.onChange(e.target.value);
        }}>
            {Object.entries(statusMap).map(([key, label]) => (
                <option key={key} value={key}>
                    {label}
                </option>
            ))}
        </select>
    )
}

export default StatusPicker;