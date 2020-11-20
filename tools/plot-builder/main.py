import matplotlib.pyplot as plt

INDEXES_NAME = [
    'Primary key on ID',
    'Index(name)',
    'Index(name, surname)',
    'Index(name, surname, id)'
]

DATA = {
    INDEXES_NAME[0]: [[1.7, 0.453, 0.587], [13, 3.39, 0.764], [13.52, 3.52, 1.4], [13.40, 3.49, 0.01]],
    INDEXES_NAME[1]: [[3.71, 0.97, 0.268], [3.69, 1.03, 1.66], [4.30, 1.12, 0.01], [4.00, 1.04, 0.01]],
    INDEXES_NAME[2]: [[40.13, 10.45, 0.024], [378.61, 98.57, 0.026], [432.19, 112.510, 0.862], [433.31, 112.81, 0.271]],
    INDEXES_NAME[3]: [[40.38, 10.51, 0.024], [381.92, 99.43, 0.026], [433.15, 112.76, 0.836], [434.40, 113.09, 0.266]],
}

PLOT_COLORS = [
    'black',
    'orange',
    'blue',
    'green'
]


def build_1_connection():
    """

    :return:
    """

    fig, axes = plt.subplots(nrows=1, ncols=3)
    fig.suptitle('Wrk results. Running 1m test with 1 thread and 1 connection')
    ax_0, ax_1, ax_2 = axes.flatten()

    width = 0.95

    ax_0.set_ylabel('Count')
    ax_0.set_title('Request/sec')
    rec_1 = ax_0.bar([0], DATA[INDEXES_NAME[0]][0][0], width, color=PLOT_COLORS[0])
    rec_2 = ax_0.bar([1], DATA[INDEXES_NAME[1]][0][0], width, color=PLOT_COLORS[1])
    rec_3 = ax_0.bar([2], DATA[INDEXES_NAME[2]][0][0], width, color=PLOT_COLORS[2])
    rec_4 = ax_0.bar([3], DATA[INDEXES_NAME[3]][0][0], width, color=PLOT_COLORS[3])
    ax_0.legend((rec_1, rec_2, rec_3, rec_4), INDEXES_NAME)

    ax_1.set_ylabel('MB')
    ax_1.set_title('Transfer/sec')
    rec_1 = ax_1.bar([0], DATA[INDEXES_NAME[0]][0][1], width, color=PLOT_COLORS[0])
    rec_2 = ax_1.bar([1], DATA[INDEXES_NAME[1]][0][1], width, color=PLOT_COLORS[1])
    rec_3 = ax_1.bar([2], DATA[INDEXES_NAME[2]][0][1], width, color=PLOT_COLORS[2])
    rec_4 = ax_1.bar([3], DATA[INDEXES_NAME[3]][0][1], width, color=PLOT_COLORS[3])
    ax_1.legend((rec_1, rec_2, rec_3, rec_4), INDEXES_NAME)

    ax_2.set_ylabel('Seconds')
    ax_2.set_title('AVG Latency')
    rec_1 = ax_2.bar([0], DATA[INDEXES_NAME[0]][0][2], width, color=PLOT_COLORS[0])
    rec_2 = ax_2.bar([1], DATA[INDEXES_NAME[1]][0][2], width, color=PLOT_COLORS[1])
    rec_3 = ax_2.bar([2], DATA[INDEXES_NAME[2]][0][2], width, color=PLOT_COLORS[2])
    rec_4 = ax_2.bar([3], DATA[INDEXES_NAME[3]][0][2], width, color=PLOT_COLORS[3])
    ax_2.legend((rec_1, rec_2, rec_3, rec_4), INDEXES_NAME)

    plt.show()


def build_10_connection():
    """

    :return:
    """

    fig, axes = plt.subplots(nrows=1, ncols=3)
    fig.suptitle('Wrk results. Running 1m test with 10 thread and 10 connection')
    ax_0, ax_1, ax_2 = axes.flatten()

    width = 0.95

    ax_0.set_ylabel('Count')
    ax_0.set_title('Request/sec')
    rec_1 = ax_0.bar([0], DATA[INDEXES_NAME[0]][1][0], width, color=PLOT_COLORS[0])
    rec_2 = ax_0.bar([1], DATA[INDEXES_NAME[1]][1][0], width, color=PLOT_COLORS[1])
    rec_3 = ax_0.bar([2], DATA[INDEXES_NAME[2]][1][0], width, color=PLOT_COLORS[2])
    rec_4 = ax_0.bar([3], DATA[INDEXES_NAME[3]][1][0], width, color=PLOT_COLORS[3])
    ax_0.legend((rec_1, rec_2, rec_3, rec_4), INDEXES_NAME)

    ax_1.set_ylabel('MB')
    ax_1.set_title('Transfer/sec')
    rec_1 = ax_1.bar([0], DATA[INDEXES_NAME[0]][1][1], width, color=PLOT_COLORS[0])
    rec_2 = ax_1.bar([1], DATA[INDEXES_NAME[1]][1][1], width, color=PLOT_COLORS[1])
    rec_3 = ax_1.bar([2], DATA[INDEXES_NAME[2]][1][1], width, color=PLOT_COLORS[2])
    rec_4 = ax_1.bar([3], DATA[INDEXES_NAME[3]][1][1], width, color=PLOT_COLORS[3])
    ax_1.legend((rec_1, rec_2, rec_3, rec_4), INDEXES_NAME)

    ax_2.set_ylabel('Seconds')
    ax_2.set_title('AVG Latency')
    rec_1 = ax_2.bar([0], DATA[INDEXES_NAME[0]][1][2], width, color=PLOT_COLORS[0])
    rec_2 = ax_2.bar([1], DATA[INDEXES_NAME[1]][1][2], width, color=PLOT_COLORS[1])
    rec_3 = ax_2.bar([2], DATA[INDEXES_NAME[2]][1][2], width, color=PLOT_COLORS[2])
    rec_4 = ax_2.bar([3], DATA[INDEXES_NAME[3]][1][2], width, color=PLOT_COLORS[3])
    ax_2.legend((rec_1, rec_2, rec_3, rec_4), INDEXES_NAME)

    plt.show()


def build_100_connection():
    """

    :return:
    """

    fig, axes = plt.subplots(nrows=1, ncols=3)
    fig.suptitle('Wrk results. Running 1m test with 10 thread and 100 connection')
    ax_0, ax_1, ax_2 = axes.flatten()

    width = 0.95

    ax_0.set_ylabel('Count')
    ax_0.set_title('Request/sec')
    rec_1 = ax_0.bar([0], DATA[INDEXES_NAME[0]][2][0], width, color=PLOT_COLORS[0])
    rec_2 = ax_0.bar([1], DATA[INDEXES_NAME[1]][2][0], width, color=PLOT_COLORS[1])
    rec_3 = ax_0.bar([2], DATA[INDEXES_NAME[2]][2][0], width, color=PLOT_COLORS[2])
    rec_4 = ax_0.bar([3], DATA[INDEXES_NAME[3]][2][0], width, color=PLOT_COLORS[3])
    ax_0.legend((rec_1, rec_2, rec_3, rec_4), INDEXES_NAME)

    ax_1.set_ylabel('MB')
    ax_1.set_title('Transfer/sec')
    rec_1 = ax_1.bar([0], DATA[INDEXES_NAME[0]][2][1], width, color=PLOT_COLORS[0])
    rec_2 = ax_1.bar([1], DATA[INDEXES_NAME[1]][2][1], width, color=PLOT_COLORS[1])
    rec_3 = ax_1.bar([2], DATA[INDEXES_NAME[2]][2][1], width, color=PLOT_COLORS[2])
    rec_4 = ax_1.bar([3], DATA[INDEXES_NAME[3]][2][1], width, color=PLOT_COLORS[3])
    ax_1.legend((rec_1, rec_2, rec_3, rec_4), INDEXES_NAME)

    ax_2.set_ylabel('Seconds')
    ax_2.set_title('AVG Latency')
    rec_1 = ax_2.bar([0], DATA[INDEXES_NAME[0]][2][2], width, color=PLOT_COLORS[0])
    rec_2 = ax_2.bar([1], DATA[INDEXES_NAME[1]][2][2], width, color=PLOT_COLORS[1])
    rec_3 = ax_2.bar([2], DATA[INDEXES_NAME[2]][2][2], width, color=PLOT_COLORS[2])
    rec_4 = ax_2.bar([3], DATA[INDEXES_NAME[3]][2][2], width, color=PLOT_COLORS[3])
    ax_2.legend((rec_1, rec_2, rec_3, rec_4), INDEXES_NAME)

    plt.show()


def build_1000_connection():
    """

    :return:
    """

    fig, axes = plt.subplots(nrows=1, ncols=3)
    fig.suptitle('Wrk results. Running 1m test with 10 thread and 1000 connection')
    ax_0, ax_1, ax_2 = axes.flatten()

    width = 0.95

    ax_0.set_ylabel('Count')
    ax_0.set_title('Request/sec')
    rec_1 = ax_0.bar([0], DATA[INDEXES_NAME[0]][3][0], width, color=PLOT_COLORS[0])
    rec_2 = ax_0.bar([1], DATA[INDEXES_NAME[1]][3][0], width, color=PLOT_COLORS[1])
    rec_3 = ax_0.bar([2], DATA[INDEXES_NAME[2]][3][0], width, color=PLOT_COLORS[2])
    rec_4 = ax_0.bar([3], DATA[INDEXES_NAME[3]][3][0], width, color=PLOT_COLORS[3])
    ax_0.legend((rec_1, rec_2, rec_3, rec_4), INDEXES_NAME)

    ax_1.set_ylabel('MB')
    ax_1.set_title('Transfer/sec')
    rec_1 = ax_1.bar([0], DATA[INDEXES_NAME[0]][3][1], width, color=PLOT_COLORS[0])
    rec_2 = ax_1.bar([1], DATA[INDEXES_NAME[1]][3][1], width, color=PLOT_COLORS[1])
    rec_3 = ax_1.bar([2], DATA[INDEXES_NAME[2]][3][1], width, color=PLOT_COLORS[2])
    rec_4 = ax_1.bar([3], DATA[INDEXES_NAME[3]][3][1], width, color=PLOT_COLORS[3])
    ax_1.legend((rec_1, rec_2, rec_3, rec_4), INDEXES_NAME)

    ax_2.set_ylabel('Seconds')
    ax_2.set_title('AVG Latency')
    rec_1 = ax_2.bar([0], DATA[INDEXES_NAME[0]][3][2], width, color=PLOT_COLORS[0])
    rec_2 = ax_2.bar([1], DATA[INDEXES_NAME[1]][3][2], width, color=PLOT_COLORS[1])
    rec_3 = ax_2.bar([2], DATA[INDEXES_NAME[2]][3][2], width, color=PLOT_COLORS[2])
    rec_4 = ax_2.bar([3], DATA[INDEXES_NAME[3]][3][2], width, color=PLOT_COLORS[3])
    ax_2.legend((rec_1, rec_2, rec_3, rec_4), INDEXES_NAME)

    plt.show()


if __name__ == '__main__':
    build_1_connection()
    build_10_connection()
    build_100_connection()
    build_1000_connection()
