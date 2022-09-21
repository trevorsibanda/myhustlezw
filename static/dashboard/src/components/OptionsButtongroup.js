

function OptionsButtonGroup(props) {
    return (
        <div class="btn-group btn-group-justified" role="group" aria-label="...">
            {props.items.map((item) => {
                return (
                    <div class="btn-group" role="group">
                        <button type="button" class={"btn "+ (props.item === (item.value ? item.value : item) ? "btn-info" : "btn-default")} onClick={() => props.onChange(item.value ? item.value : item)}>{item.component ? item.component : item}</button>
                    </div>
                )
            })}
        </div>
    )
}

export default OptionsButtonGroup;