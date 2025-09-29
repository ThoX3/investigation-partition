def format_size(size_bytes):
    if size_bytes < 1024:
        return f"{size_bytes} B"
    elif size_bytes < 1048576:
        return f"{size_bytes/1024:.1f} KB"
    elif size_bytes < 1073741824:
        return f"{size_bytes/1048576:.1f} MB"
    else:
        return f"{size_bytes/1073741824:.1f} GB"
