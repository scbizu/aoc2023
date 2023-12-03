const std = @import("std");

const IndexNumber = struct {
    index: u32,
    number: u8,
};

const Number = struct {
    numberStr: []u8,
    number: u8,
};

fn indexNumberAscending(lhs: IndexNumber, rhs: IndexNumber) bool {
    return lhs.index < rhs.index;
}

pub fn main() !void {
    var file = try std.fs.cwd().openFile("input.txt", .{});
    defer file.close();
    var reader = std.io.bufferedReader(file.reader());
    var stream = reader.reader();
    var buffer: [1024]u8 = undefined;
    var total: u32 = 0;
    const numberSet = std.ArrayList(Number).init(std.heap.page_allocator);
    defer numberSet.deinit();
    numberSet.insert(.{ .numberStr = "one", .number = 1 });
    numberSet.insert(.{ .numberStr = "two", .number = 2 });
    numberSet.insert(.{ .numberStr = "three", .number = 3 });
    numberSet.insert(.{ .numberStr = "four", .number = 4 });
    numberSet.insert(.{ .numberStr = "five", .number = 5 });
    numberSet.insert(.{ .numberStr = "six", .number = 6 });
    numberSet.insert(.{ .numberStr = "seven", .number = 7 });
    numberSet.insert(.{ .numberStr = "eight", .number = 8 });
    numberSet.insert(.{ .numberStr = "nine", .number = 9 });
    while (try stream.readUntilDelimiterOrEof(&buffer, '\n')) |line| {
        var array = std.ArrayList(IndexNumber).init(std.heap.page_allocator);
        defer array.deinit();
        for (line, 0..) |c, i| {
            if (c >= '0' and c <= '9') {
                try array.append(.{ .index = i, .number = c - '0' });
            }
        }
        for (numberSet.items) |number| {
            // if line contains number
            const index = try std.mem.indexOf([]u8, line, number.numberStr);
            if (index != null) {
                try array.append(.{ .index = index, .number = number.number });
            }
        }

        if (array.items.len == 1) {
            try array.append(array.items[0]);
        }
        std.sort.insertion(IndexNumber, array.items, {}, indexNumberAscending);
        const str = [_]u8{ array.items[0], array.items[1] };
        const num = try std.fmt.parseInt(u32, &str, 10);
        total += num;
    }
    std.debug.print("total: {d}\n", .{total});
}
