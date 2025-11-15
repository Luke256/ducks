import { PosterStatusLabels } from "@/types/poster";

type Props = {
    status: string;
    onChange: (newStatus: string) => void;
}

const StatusPicker = (props: Props) => {

    return (
        <select value={props.status} className="border border-gray-300 p-2 hover:cursor-pointer" onChange={(e) => {
            props.onChange(e.target.value);
        }}>
            {Object.entries(PosterStatusLabels).map(([key, label]) => (
                <option key={key} value={key}>
                    {label}
                </option>
            ))}
        </select>
    )
}

export default StatusPicker;