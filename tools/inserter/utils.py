

def divide_sequence_into_chunks(sequence: list, size: int):
    for i in range(0, len(sequence), size):
        yield sequence[i:i + size]