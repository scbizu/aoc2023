const std = @import("std");

pub fn main() !void {
    var file = try std.fs.cwd().openFile("input.txt", .{});
    defer file.close();
    var reader = std.io.bufferedReader(file.reader());
    var stream = reader.reader();
    var buffer: [1024]u8 = undefined;
    var total: u32 = 0;
    while (try stream.readUntilDelimiterOrEof(&buffer, '\n')) |line| {
        var array = std.ArrayList(u8).init(std.heap.page_allocator);
        defer array.deinit();
        for (line) |c| {
            if (c >= '0' and c <= '9') {
                try array.append(c);
            }
        }
        if (array.items.len == 1) {
            try array.append(array.items[0]);
        }
        const str = [_]u8{ array.items[0], array.items[1] };
        const num = try std.fmt.parseInt(u32, &str, 10);
        total += num;
    }
    std.debug.print("total: {d}\n", .{total});
}
